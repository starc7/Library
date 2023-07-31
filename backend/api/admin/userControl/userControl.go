package userControl

import (
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"champ.com/assets"
)


// GET All Users Data by Admin
func GetAllUserByAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []assets.User
		adminID := c.Query("adminID")

		rows, err := db.Query("SELECT * FROM Users WHERE LibID = $1", adminID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var user assets.User

			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.ContactNumber, &user.Role, &user.LibID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}
		c.JSON(http.StatusOK, users)
	}
}


// POST Add New User by Admin
func AddNewUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		Name := c.PostForm("Name")
		Email := c.PostForm("Email")
		ContactNumber := c.PostForm("ContactNumber")
		Role := c.PostForm("Role")
		LibID := c.PostForm("LibID")
		
		var exists bool

		err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM Library WHERE ID = $1)", LibID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		if exists {
			_, err = db.Exec("INSERT INTO Users(Name, Email, ContactNumber, Role, LibID) VALUES ($1, $2, $3, $4, $5)", 
					Name, Email, ContactNumber, Role, LibID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}
			
			c.JSON(http.StatusOK, gin.H{"message": "New User added successfully"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Librarian ID"})
		}
	}
}



