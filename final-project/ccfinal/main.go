package main

import (
	"fmt"
	"httpmonitor/model"
	"httpmonitor/pkg/conf"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	mdl model.Model
	cfg conf.Config
)

func checkURL(address string) bool {
	url := "http://" + address
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func monitor(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		now := time.Now()
		servers, err := mdl.GetAll()
		if err != nil {
			log.Println("Could not recieve servers from DB")
		}
		log.Println("Checking URLs at", now.Format("15:04:05"))
		for _, server := range servers {
			result := checkURL(server.Address)
			log.Println(server.Address, ":", result)
			if result {
				mdl.SubmitSuccess(server)
			} else {
				mdl.SubmitFailure(server)
			}
		}
		// Wait for the next tick
		<-t.C
	}
}

func getServerByIDHandler(c echo.Context) error {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		log.Println("bad ID input")
		return c.JSON(http.StatusBadRequest, "ID should be all numbers.")
	}
	server, err := mdl.GetByID(uint(id))
	if err != nil {
		log.Println("No servert was submitted with this ID")
		return c.JSON(http.StatusNotFound, "No servert was submitted with this ID")
	} else {
		log.Println("returning server" + fmt.Sprint(id) + "to the client.")
		return c.JSON(http.StatusFound, server)
	}
}

func getAllServersHandler(c echo.Context) error {
	servers, err := mdl.GetAll()
	if err != nil {
		log.Println("Could not get all servers from db")
		return c.JSON(http.StatusInternalServerError, "Error connecting to db")
	} else {
		return c.JSON(http.StatusOK, servers)
	}
}

func createServerHandler(c echo.Context) error {
	address := c.FormValue("address")
	id, err := mdl.CreateServer(address)
	if err != nil {
		log.Println("Could not create the server")
		return c.JSON(http.StatusInternalServerError, "Error connecting to db")
	} else {
		log.Println("server created with id:" + fmt.Sprint(id))
		return c.JSON(http.StatusOK, "server with id:"+fmt.Sprint(id)+" was created successfully.")
	}

}

func main() {
	var err error
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	cfg = conf.Load()
	for {
		mdl, err = model.Initialize(cfg, dbUser, dbPassword, dbName)
		if err != nil {
			log.Println("could not connect to db")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	go monitor(cfg.Monitor.IntervalM)
	e := echo.New()
	e.GET("/api/server/:id", getServerByIDHandler)
	e.GET("/api/server/all", getAllServersHandler)
	e.POST("/api/server/", createServerHandler)
	e.Start("0.0.0.0:" + fmt.Sprint(cfg.Monitor.Port))
}
