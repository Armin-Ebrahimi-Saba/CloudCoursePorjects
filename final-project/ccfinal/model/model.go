package model

import (
	"errors"
	"fmt"
	"httpmonitor/pkg/conf"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Model struct {
	writeDB *gorm.DB
	readDB  *gorm.DB
}

type Server struct {
	gorm.Model
	Address     string
	Success     int
	Failure     int
	LastFailure time.Time
}

func Initialize(cfg conf.Config, user string, pass string, dbname string) (Model, error) {

	var model Model
	rdsn := user + ":" + pass + "@tcp(" + cfg.Mysql.RAddr + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	rdb, err := gorm.Open(mysql.Open(rdsn))
	if err != nil {
		return model, err
	}
	rdb.AutoMigrate(&Server{})
	wdsn := user + ":" + pass + "@tcp(" + cfg.Mysql.WAddr + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	wdb, err := gorm.Open(mysql.Open(wdsn))
	if err != nil {
		return model, err
	}
	wdb.AutoMigrate(&Server{})
	model.readDB = rdb
	model.writeDB = wdb
	return model, nil
}

func (model Model) CreateServer(address string) (uint, error) {
	t := time.Date(2000, 1, 1, 1, 1, 1, 1, time.UTC)
	server := Server{Address: address, Success: 0, Failure: 0, LastFailure: t}
	result := model.writeDB.Create(&server)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return server.ID, nil
}

func (model Model) GetByID(id uint) (Server, error) {
	var server Server
	result := model.readDB.Where("ID = ?", fmt.Sprint(id)).First(&server)
	if server.ID != id {
		log.Println(result.Error, "no user with this ID")
		return server, errors.New("no user with this ID")
	}
	return server, nil
}

func (model Model) GetAll() ([]Server, error) {
	var servers []Server
	result := model.readDB.Find(&servers) // SELECT * FROM servers;
	if result.Error != nil {
		return nil, result.Error
	}
	return servers, nil
}

func (model Model) SubmitFailure(server Server) error {
	server.Failure += 1
	server.LastFailure = time.Now()
	result := model.writeDB.Save(&server)
	return result.Error
}

func (model Model) SubmitSuccess(server Server) error {
	server.Success += 1
	result := model.writeDB.Save(&server)
	return result.Error
}
