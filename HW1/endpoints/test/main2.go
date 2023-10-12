

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2
var x = 0

func CreateShortURL(c *gin.Context) {
	log.Printf("Got a %vth request.", x)
	x += 1
	c.JSON(http.StatusOK, gin.H{
		"message": "Hi",
	})
}

func init() {
	r := gin.Default()

	r.GET("/app", CreateShortURL)
	// r.GET("/app/:shortcode", GetShortURL)
	ginLambda = ginadapter.NewV2(r)
}

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}

// func main() {
// 	router := gin.Default()

// 	router.Group("/").GET("login", func(ctx *gin.Context) {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"message": "/login",
// 		})
// 	})
// 	// v := router.Group("/")
// 	// {
// 	// 	v.GET("/login", func(ctx *gin.Context) {
// 	// 		ctx.JSON(http.StatusOK, gin.H{
// 	// 			"message": "/login",
// 	// 		})
// 	// 	})
// 	// 	v.GET("/read", func(ctx *gin.Context) {
// 	// 		ctx.JSON(http.StatusOK, gin.H{
// 	// 			"message": "/read",
// 	// 		})
// 	// 	})
// 	// }

// 	v1 := router.Group("/v1")
// 	{
// 		v1.GET("/login", func(ctx *gin.Context) {
// 			ctx.JSON(http.StatusOK, gin.H{
// 				"message": "v1/login",
// 			})
// 		})
// 		v1.GET("/read", func(ctx *gin.Context) {
// 			ctx.JSON(http.StatusOK, gin.H{
// 				"message": "v1/read",
// 			})
// 		})
// 	}

// 	v2 := router.Group("/v2")
// 	{
// 		v2.GET("/login", func(ctx *gin.Context) {
// 			ctx.JSON(http.StatusOK, gin.H{
// 				"message": "v2/login",
// 			})
// 		})
// 		v2.GET("/read", func(ctx *gin.Context) {
// 			ctx.JSON(http.StatusOK, gin.H{
// 				"message": "v2/read",
// 			})
// 		})
// 	}

// 	router.Run()
// }
