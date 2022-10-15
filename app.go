package main

import (
	"context"
	"io/ioutil"
	"net/http" // for returning standard defined api response codes
	"time"

	flagsmith "github.com/Flagsmith/flagsmith-go-client/v2"
	"github.com/gin-gonic/gin" // to easily bootstrap api server
	"gopkg.in/yaml.v3"
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

func loadConfigMap() (map[string]string, error) {
	yfile, err := ioutil.ReadFile("config.yml")

	if err != nil {
		return nil, err
	}

	config := make(map[string]string)

	err = yaml.Unmarshal(yfile, &config)

	if err != nil {
		return nil, err
	}
	return config, nil
}

// getBooks responds with the list of all books as JSON.
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books) // IndentedJSON is like pretty-print which makes it more readable
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := loadConfigMap()
	if err != nil {
		panic("could not read config")
	}

	// Initialise the Flagsmith client
	client := flagsmith.NewClient(config["key"],
		flagsmith.WithLocalEvaluation(), // for local evaluation
		flagsmith.WithEnvironmentRefreshInterval(30*time.Second),
		flagsmith.WithContext(ctx))

	router := gin.Default()                      // creates a gin engine instance and returns it
	router.GET("/books", getBooks)               // registering the function "getBooks" that will be called when /books endpoint is hit
	router.POST("/books", func(c *gin.Context) { // defining the function to be executed when /books POST endpoint is hit
		flags, err := client.GetEnvironmentFlags()
		if err != nil {
			return
		}
		isEnabled, err := flags.IsFeatureEnabled("enable_new_books")
		if err != nil {
			return
		}
		// Add the new book to the list if feature flag is enabled
		if isEnabled {
			var newBook book
			if err := c.BindJSON(&newBook); err != nil {
				return
			}
			books = append(books, newBook)
			c.IndentedJSON(http.StatusCreated, newBook)
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"code":    http.StatusMethodNotAllowed,
				"message": "sorry, please come back later",
			})
		}
	})

	router.Run("localhost:8080") // running the server on port 8080

}
