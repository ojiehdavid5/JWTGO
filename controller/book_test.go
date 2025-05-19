package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/chuks/JWTGO/controller"
	"github.com/chuks/JWTGO/model"
	// "github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)


func TestBookController(t *testing.T) {
	// Setup test database and Fiber app
	db := setupTestDB(t)
	app := setupFiberApp(db)

	// bookController := controller.NewBook(db)

	// Test CreateBook
	t.Run("CreateBook", func(t *testing.T) {
		reqBody := `{"title":"Test Book", "author":"Test Author"}`
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Check if the book was created in the database
		var createdBook model.Book
		db.First(&createdBook, "title = ?", "Test Book")
		assert.Equal(t, "Test Book", createdBook.Title)
		assert.Equal(t, "Test Author", createdBook.Author)
	})

	// Test GetBooks
	t.Run("GetBooks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var books []model.Book
		err = json.NewDecoder(resp.Body).Decode(&books)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(books))
		assert.Equal(t, "Test Book", books[0].Title)
	})

	// Test GetBook by ID
	t.Run("GetBookByID", func(t *testing.T) {
		var createdBook model.Book
		db.First(&createdBook, "title = ?", "Test Book")

		req := httptest.NewRequest(http.MethodGet, "/books/"+strconv.Itoa(int(createdBook.ID)), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var responseBook model.Book
		err = json.NewDecoder(resp.Body).Decode(&responseBook)
		assert.NoError(t, err)
		assert.Equal(t, createdBook.Title, responseBook.Title)
		assert.Equal(t, createdBook.Author, responseBook.Author)
	})

	// Test UpdateBook
	t.Run("UpdateBook", func(t *testing.T) {
		var createdBook model.Book
		db.First(&createdBook, "title = ?", "Test Book")

		reqBody := `{"title":"Updated Book", "author":"Updated Author"}`
		req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(int(createdBook.ID)), bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var updatedBook model.Book
		db.First(&updatedBook, createdBook.ID)
		assert.Equal(t, "Updated Book", updatedBook.Title)
		assert.Equal(t, "Updated Author", updatedBook.Author)
	})

	// Test DeleteBook
	t.Run("DeleteBook", func(t *testing.T) {
		var createdBook model.Book
		db.First(&createdBook, "title = ?", "Updated Book")

		req := httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(int(createdBook.ID)), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check if the book was deleted from the database
		var deletedBook model.Book
		res := db.First(&deletedBook, createdBook.ID)
		assert.Error(t, res.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, res.Error)
	})
}