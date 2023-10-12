package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Form struct {
	Email      string                `form:"email" binding:"required,email"`
	LastName   string                `form:"lastName" binding:"required,min=5,max=15"`
	NationalID int                   `form:"nationalID" binding:"required,min=0,max=100000"`
	Image1     *multipart.FileHeader `form:"image1" binding:"required"`
	Image2     *multipart.FileHeader `form:"image2" binding:"required"`
}

type User struct {
	Username   string `dynamodbav:"username"`
	Email      string `dynamodbav:"email"`
	LastName   string `dynamodbav:"lastName"`
	NationalID int    `dynamodbav:"nationalID"`
	IP         string `dynamodbav:"ip"`
	Image1     string `dynamodbav:"image1"`
	Image2     string `dynamodbav:"image2"`
	State      string `dynamodbav:"state"`
}

var initialized = false

const TableName = "BankingAuthenticationService"

var bucketName = "banking-authentication-images"
var QueueName = "banking-authentication"

var ginLambda *ginadapter.GinLambdaV2
var db dynamodb.Client
var s3Client *s3.Client
var sqsClient *sqs.Client

func main() {
	lambda.Start(Handler)
}

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	db = *dynamodb.NewFromConfig(sdkConfig)
	s3Client = s3.NewFromConfig(sdkConfig)
	sqsClient = sqs.NewFromConfig(sdkConfig)
}

func SaveImage(c *gin.Context) (*User, error) {
	var form Form
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("Couldn't open file Here's why: %v\n", err.Error())
		return nil, err
	}
	f, err := form.Image1.Open()
	if err != nil {
		log.Printf("Couldn't open file Here's why: %v\n", err.Error())
		return nil, err
	}
	newUsername := fmt.Sprint(form.NationalID) // c.RemoteIP() + time.Now().UTC().String()
	keyName := newUsername + "1"
	var partMiBs int64 = 10
	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName + filepath.Ext(form.Image1.Filename)),
		Body:   f,
	})
	if err != nil {
		log.Printf("Couldn't upload file to %v Here's why: %v\n",
			bucketName, err)
		return nil, err
	}
	url1 := result.Location
	keyName = newUsername + "2"
	f, err = form.Image2.Open()
	if err != nil {
		log.Printf("Couldn't open file Here's why: %v\n", err.Error())
		return nil, err
	}
	result, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName + filepath.Ext(form.Image1.Filename)),
		Body:   f,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v. Here's why: %v\n",
			f, bucketName, err)
		return nil, err
	}
	url2 := result.Location
	user := User{
		Username:   newUsername,
		Email:      form.Email,
		LastName:   form.LastName,
		NationalID: form.NationalID,
		IP:         c.RemoteIP(),
		Image1:     url1,
		Image2:     url2,
		State:      "Pending",
	}
	return &user, nil
}

func InsertRecord(user *User, c *gin.Context) error {
	item, err := attributevalue.MarshalMap(*user)
	if err != nil {
		log.Printf("Couldn't marshal map to item. Here's why: %v\n", err)
		return err
	}
	_, err = db.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	if err != nil {
		log.Printf("Couldn't insert record. Here's why: %v\n", err)
		return err
	}
	return nil
}

func EnqueueRequests(c *gin.Context) {
	user, err := SaveImage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = InsertRecord(user, c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = Push(user.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Your request was registered.",
	})
}

func GetStatus(c *gin.Context) {
	id := c.Query("nationalID")
	username, err := attributevalue.Marshal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		log.Printf("Couldn't marshal nationalID. Here's why: %v\n", err)
		return
	}
	result, err := db.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]dynamotypes.AttributeValue{
			"username": username,
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		log.Printf("Couldn't retrieve record from dynamodb. Here's why: %v\n", err)
		return
	}
	var user User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Printf("Couldn't unmarshal resposne from dynamodb. Here's why: %v\n", err)
		return
	}
	if c.RemoteIP() != user.IP {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid IP address.",
		})
		log.Printf("No match between requester's IP and record's IP")
		return
	}
	var message string
	switch user.State {
	case "Pending":
		message = "Your request is pending."
	case "Approved":
		message = "Your request was approved."
	case "Rejected":
		message = "Your request was rejected. Please try again later."
	default:
		message = "Please retry later."
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
}

func Handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ginEngine := BuildEngine()
	if !initialized {
		ginEngine.SetTrustedProxies(nil)
		ginEngine.POST("/register", EnqueueRequests)
		ginEngine.GET("/register", GetStatus)
		ginLambda = ginadapter.NewV2(ginEngine)
		initialized = true
	}
	return ginLambda.ProxyWithContext(ctx, request)
}

func BuildEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	return engine
}

func Push(username string) error {
	// Get URL of queue
	gQInput := &sqs.GetQueueUrlInput{
		QueueName:              &QueueName,
		QueueOwnerAWSAccountId: aws.String("106719423561"),
	}
	result, err := sqsClient.GetQueueUrl(context.TODO(), gQInput)
	if err != nil {
		return err
	}
	queueURL := result.QueueUrl

	sMInput := &sqs.SendMessageInput{
		DelaySeconds: 6,
		MessageAttributes: map[string]sqstypes.MessageAttributeValue{
			"Username": {
				DataType:    aws.String("String"),
				StringValue: aws.String(username),
			},
		},
		MessageBody: aws.String(username),
		QueueUrl:    queueURL,
	}

	_, err = sqsClient.SendMessage(context.TODO(), sMInput)
	if err != nil {
		return err
	}
	return nil
}

// fmt.Println("Sent message with ID: " + *resp.MessageId)

// type SQSSendMessageAPI interface {
// 	GetQueueUrl(ctx context.Context,
// 		params *sqs.GetQueueUrlInput,
// 		optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

// 	SendMessage(ctx context.Context,
// 		params *sqs.SendMessageInput,
// 		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
// }

// func GetQueueURL(c context.Context, api SQSSendMessageAPI, input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
// 	return api.GetQueueUrl(c, input)
// }

// func SendMsg(c context.Context, api SQSSendMessageAPI, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
// 	return api.SendMessage(c, input)
// }
