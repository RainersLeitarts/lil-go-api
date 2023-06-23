package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Title one!", Author: "dunno", Quantity: 3},
	{ID: "2", Title: "Can you even read?", Author: "thinkaboutit", Quantity: 2},
	{ID: "3", Title: "Boooooooooooooook", Author: "writer", Quantity: 5},
}

func findBookById(id string) (*book, error) {
	for index, book := range books {
		if book.ID == id {
			return &books[index], nil
		}
	}

	return nil, errors.New("Book not found")
}

func getBooks(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, books)
}

func getBook(context *gin.Context) {
	id := context.Param("id")
	book, err := findBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, *book)
}

func setBook(context *gin.Context) {
	var newBook book

	if err := context.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	context.IndentedJSON(http.StatusCreated, newBook)
}

func checkoutBook(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Query parameter 'id' not found"})
		return
	}

	book, err := findBookById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book out of stock"})
		return
	}

	book.Quantity--

	context.IndentedJSON(http.StatusOK, book)
}

func main() {
	baseRoute := "localhost"
	port := "8080"
	router := gin.Default()

	// Routes
	router.GET("/books", getBooks)
	router.GET("/book/:id", getBook)
	router.POST("/book", setBook)
	router.PATCH("/checkout", checkoutBook)
	router.Run(baseRoute + ":" + port)
}
