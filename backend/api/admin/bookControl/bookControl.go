package bookControl

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"champ.com/assets"
)


// GET List of All books available in Library
func GetAllBooksByAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")
		var books []assets.BookInventory

		rows, err := db.Query("SELECT * FROM BookInventory WHERE LibID = $1", adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var book assets.BookInventory

			if err := rows.Scan(&book.ISBN, &book.LibID, &book.Title, &book.Authors, &book.Publisher, &book.Version, &book.TotalCopies, &book.AvailableCopies); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			books = append(books, book)
		}
		c.JSON(http.StatusOK, books)
	}
}


// POST Add new book in Library
func AddNewBookToLib(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ISBN := c.PostForm("ISBN")
		Title := c.PostForm("Title")
		Authors := c.PostForm("Authors")
		Publisher := c.PostForm("Publisher")
		Version := c.PostForm("Version")
		TotalCopies := c.PostForm("TotalCopies")
		AvailableCopies := c.PostForm("AvailableCopies")
		LibID := c.PostForm("LibID")

		var libExists bool
		var bookExists bool

		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Library WHERE ID = $1)", LibID).Scan(&libExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM BookInventory WHERE ISBN = $1)", ISBN).Scan(&bookExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		

		if libExists && !bookExists {
			_, err = db.Exec("INSERT INTO BookInventory (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
				ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Book added successfully to the inventory"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Library with the specified ID does not exist or Book already exists"})
		}
	}
}


// PUT Update the Record of Books
func UpdateBookRecord(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ISBN := c.Query("ISBN")
		Copies := c.Query("Copies")

		var bookExists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM BookInventory WHERE ISBN = $1)", ISBN).Scan(&bookExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		if bookExists {
			_, err := db.Exec("UPDATE BookInventory SET AvailableCopies = AvailableCopies + $1, TotalCopies = TotalCopies + $2 WHERE ISBN = $3", Copies, Copies, ISBN)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Book Updated Succesfully"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid ISBN"})
		}
	}
}


// PUT Book Return By ID
func BookReturnByUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Query("issueID")
		returnDate := c.Query("returnDate")
		returnReaderID := c.Query("returnReaderID")
		issueStatus := "Submitted"

		var issueIDExists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM IssueRegistry WHERE IssueID = $1)", issueID).Scan(&issueIDExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if issueIDExists {
			_, err = db.Exec("UPDATE IssueRegistry SET ReturnDate = $1, ReturnApproverID = $2, IssueStatus = $3 WHERE IssueID = $4", returnDate, returnReaderID, issueStatus, issueID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			_, err = db.Exec("UPDATE BookInventory SET TotalCopies = TotalCopies+1, AvailableCopies = AvailableCopies+1 WHERE ISBN = (SELECT BookID FROM IssueRegistry WHERE IssueID = $1 LIMIT 1)", issueID)
			// UPDATE BookInventory BI SET BI.TotalCopies = BI.TotalCopies+1, BI.AvailableCopies = BI.AvailableCopies+1 FROM IssueRegistry IR WHERE IR.BookID = BI.ISBN AND IR.IssueID = $1

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}


			c.JSON(http.StatusOK, gin.H{"message": "Book returned successfully"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error" : "Invalid Issue ID"})
		}

	}
}
