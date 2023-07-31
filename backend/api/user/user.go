package user

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"github.com/agnivade/levenshtein"
	"strings"
	"champ.com/assets"
)


// GET Search Book by Author Name
func SearchBookByAuthorsName(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []assets.BookInventory

		authorName := c.Query("authors")

		rows, err := db.Query("SELECT * FROM BookInventory")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var book assets.BookInventory

			if err := rows.Scan(&book.ISBN, &book.LibID, &book.Title, &book.Authors, &book.Publisher, 
				&book.Version, &book.TotalCopies, &book.AvailableCopies); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			matchScore := levenshtein.ComputeDistance(strings.ToLower(authorName), strings.ToLower(book.Authors))
			if matchScore <= 2 {
				books = append(books, book)
			}
		}
		c.JSON(http.StatusOK, books)
	}
}


// GET Search Book By Title
func SearchBookByTitle(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []assets.BookInventory

		title := c.Query("title")

		rows, err := db.Query("SELECT * FROM BookInventory")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var book assets.BookInventory

			if err := rows.Scan(&book.ISBN, &book.LibID, &book.Title, &book.Authors, &book.Publisher, 
				&book.Version, &book.TotalCopies, &book.AvailableCopies); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			matchScore := levenshtein.ComputeDistance(strings.ToLower(title), strings.ToLower(book.Title))
			if matchScore <= 2 {
				books = append(books, book)
			}
		}
		c.JSON(http.StatusOK, books)
	}
}


// GET Search Books by Publisher
func SearchBookByPublisher(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []assets.BookInventory

		publisher := c.Query("publisher")

		rows, err := db.Query("SELECT * FROM BookInventory")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var book assets.BookInventory

			if err := rows.Scan(&book.ISBN, &book.LibID, &book.Title, &book.Authors, &book.Publisher, 
				&book.Version, &book.TotalCopies, &book.AvailableCopies); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			matchScore := levenshtein.ComputeDistance(strings.ToLower(publisher), strings.ToLower(book.Publisher))
			if matchScore <= 2 {
				books = append(books, book)
			}
		}
		c.JSON(http.StatusOK, books)
	}
}


// POST Raise request for Book
func RaiseRequestForBook(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		BookID := c.PostForm("BookID")
		ReaderID := c.PostForm("ReaderID")
		RequestDate := c.PostForm("RequestDate")
		RequestType := c.PostForm("RequestType")

		var readerIDexists bool
		var isBookAvailable bool

		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Users WHERE ID = $1)", ReaderID).Scan(&readerIDexists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM BookInventory WHERE ISBN = $1 AND AvailableCopies > 0)", BookID).Scan(&isBookAvailable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		if readerIDexists {
			if isBookAvailable {
				_, err = db.Exec("INSERT INTO RequestEvents (BookId, ReaderID, RequestDate, RequestType) VALUES($1, $2, $3, $4)", 
				BookID, ReaderID, RequestDate, RequestType)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"Message": "Request raised successfully"})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Message": "Book is not available or you are looking for wrong Book"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Reader/UserID"})
		}
	}
}


// GET Issue Registry for user
func GetIssueRegistryByUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("userID")

		rows, err := db.Query("SELECT IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate, ExpectedReturnDate FROM IssueRegistry WHERE ReturnDate IS NULL AND ReaderID = $1", userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		defer rows.Close()

		var issues []assets.PendingBooks
		for rows.Next() {
			var issue assets.PendingBooks
			
			if err := rows.Scan(&issue.IssueID, &issue.ISBN, &issue.ReaderID, &issue.IssueApproverID, &issue.IssueStatus, &issue.IssueDate, 
				&issue.ExpectedReturnDate); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}

			issue.IssueDate = issue.IssueDate[0:10]
			issue.ExpectedReturnDate = issue.ExpectedReturnDate[0:10]
			
			issues = append(issues, issue)
		}
		c.JSON(http.StatusOK, issues)
	}
}
