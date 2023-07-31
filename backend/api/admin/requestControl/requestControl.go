package requestControl


import (
	"fmt"
	//"time"
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"champ.com/assets"
)


// GET all pending requests
func GetAllPendingRequest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")

		rows, err := db.Query("SELECT ReqID, BookID, ReaderID, RequestDate, RequestType FROM RequestEvents WHERE ApproverID IS NULL AND ReaderID IN (SELECT ID FROM Users WHERE LibID = $1)", 
		adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var requests []assets.PendingRequest
		for rows.Next() {
			var request assets.PendingRequest
			if err := rows.Scan(&request.ReqID, &request.BookID, &request.ReaderID, &request.RequestDate, &request.RequestType); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			request.RequestDate = request.RequestDate[0:10]
			requests = append(requests, request)
		}
		c.JSON(http.StatusOK, requests)
	}
}


// GET List of all request
func GetAllRequests(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")

		rows, err := db.Query("SELECT * FROM RequestEvents WHERE ApproverID IS NOT NULL AND ReaderID IN (SELECT ID FROM Users WHERE LibID = $1)", 
		adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var requests []assets.RequestEvents
		for rows.Next() {
			var request assets.RequestEvents
			if err := rows.Scan(&request.ReqID, &request.BookID, &request.ReaderID, &request.RequestDate, &request.ApprovalDate, &request.ApproverID ,&request.RequestType); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			request.RequestDate = request.RequestDate[0:10]
			request.ApprovalDate = request.ApprovalDate[0:10]

			requests = append(requests, request)
		}
		c.JSON(http.StatusOK, requests)
	}
}


// PUT Approve Request by Admin
func ApproveRequstByAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// AdminID := c.PostForm("AdminID")
		RequestID := c.PostForm("RequestID")
		IssueDate := c.PostForm("IssueDate")
		ExpectedReturnDate := c.PostForm("ExpectedReturnDate")

		var reqIDexists bool
		var ISBN string
		var ReaderID int

		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM RequestEvents WHERE ReqID = $1)", RequestID).Scan(&reqIDexists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if reqIDexists {
			err = db.QueryRow("SELECT BookId, ReaderID FROM RequestEvents WHERE ReqID = $1", RequestID).Scan(&ISBN, &ReaderID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_,err = db.Exec("INSERT INTO IssueRegistry(ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate, ExpectedReturnDate) VALUES($1, $2, $3, 'Near Reader', $4, $5)", 
			ISBN, ReaderID, ReaderID, IssueDate, ExpectedReturnDate)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_, err = db.Exec("UPDATE BookInventory SET TotalCopies = TotalCopies-1, AvailableCopies = AvailableCopies-1 WHERE ISBN = $1", ISBN)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			_, err = db.Exec("UPDATE RequestEvents SET ApprovalDate = $1, ApproverID = $2 WHERE ReqID = $3", IssueDate, ReaderID, RequestID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Request ID %s has been approved successfully and Book ID %s has been issued to the User %d", RequestID, ISBN, ReaderID)})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Request ID"})
		}
	}
}


// DELETE Request
func RejectRequest(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")
		requestID := c.Query("requestID")

		var adminExists bool
		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Library WHERE ID = $1)", adminID).Scan(&adminExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if adminExists {
			_, err := db.Exec("DELETE FROM RequestEvents WHERE ReqID = $1", requestID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Request rejected"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Admin ID"})
		}
	}
}

