package main

import (
	"net/http" // for returning standard defined api response codes
	"github.com/gin-gonic/gin" // to easily bootstrap api server
)

// a book struct(class) which contains attributes describing a book
type book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

// a list with 3 distinct books defined as per the struct above
var books = []book{
	{ID: "1", Title: "Harry Potter", Author: "J. K. Rowling", Price: 26.99},
	{ID: "2", Title: "War and Peace", Author: "Leo Tolstoy", Price: 17.99},
	{ID: "3", Title: "The Kite Runner", Author: "Khaled Hosseini", Price: 29.99},
}

// getBooks responds with the list of all books as JSON.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books) // IndentedJSON is like pretty-print which makes it more readable
}

func main() {
	router := gin.Default()        // creates a gin engine instance and returns it
	router.GET("/books", getBooks) // registering the function "getBooks" that will be called when /books endpoint is hit
	router.Run("localhost:8080")   // running the server on port 8080
}