package issueControl

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"champ.com/assets"
)


// GET all details of Registry
func GetAllIssueRegistry(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")

		rows, err := db.Query("SELECT * FROM IssueRegistry WHERE ReturnDate IS NOT NULL AND ReaderID IN (SELET ID FROM Users WHERE LibID = $1)", adminID)
		// SELECT IR.* FROM IssueRegistry IR JOIN Users U ON IR.ReaderID = U.ID WHERE IR.ReturnDate IS NOT NULL AND U.LibID = $1, adminID
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		defer rows.Close()

		var issues []assets.IssueRegistry
		for rows.Next() {
			var issue assets.IssueRegistry
			
			if err := rows.Scan(&issue.IssueID, &issue.ISBN, &issue.ReaderID, &issue.IssueApproverID, &issue.IssueStatus, &issue.IssueDate, 
				&issue.ExpectedReturnDate, &issue.ReturnDate, &issue.ReturnApproverID); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
			
			issue.IssueDate = issue.IssueDate[0:10]
			issue.ExpectedReturnDate = issue.ExpectedReturnDate[0:10]
			issue.ReturnDate = issue.ReturnDate[0:10]

			issues = append(issues, issue)
		}
		c.JSON(http.StatusOK, issues)
	}
}


// GET Pending Books Details
func GetAllPendingBooksDetails(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")

		rows, err := db.Query("SELECT IssueID, ISBN, ReaderID, IssueApproverID, IssueStatus, IssueDate, ExpectedReturnDate FROM IssueRegistry WHERE ReturnDate IS NULL AND ReaderID IN (SELET ID FROM Users WHERE LibID = $1)", adminID)
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


