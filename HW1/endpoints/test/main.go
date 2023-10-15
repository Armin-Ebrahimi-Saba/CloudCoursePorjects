package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	req := c.Request
	// Print the request method (GET, POST, etc.)
	log.Println("Request IP:", c.RemoteIP())
	// Print the request method (GET, POST, etc.)
	log.Println("Request Method:", req.Method)

	// Print the request URL
	log.Println("Request URL:", req.URL.String())

	// Print the request headers
	log.Println("Request Headers:")
	for key, values := range req.Header {
		log.Println(key, ":", values)
	}

	// Print the query parameters
	log.Println("Query Parameters:")
	query := req.URL.Query()
	for key, values := range query {
		log.Println(key, ":", values)
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.Run("localhost:8080")

}
