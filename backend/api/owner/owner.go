package owner

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"database/sql"
	"net/http"
	"champ.com/assets"
)


// GET of all Admin Data
func GetAllAdminData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var admins []assets.Library

		rows, err := db.Query("SELECT * FROM Library")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		
		for rows.Next() {
			var admin assets.Library

			if err := rows.Scan(&admin.ID, &admin.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			
			admins = append(admins, admin)
		}
		c.JSON(http.StatusOK, admins)
	}
}


// GET all Users Data
func GetAllUserData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []assets.User

		rows, err := db.Query("SELECT * FROM Users")
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


// GET All Books data
func GetAllBookData(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []assets.BookInventory

		rows, err := db.Query("SELECT * FROM BookInventory")
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


// POST Add New Admin by using his Name only
func AddNewAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("password")

		if password == "OwnerPassword" {
			adminName := c.Query("name")

			_, err := db.Exec("INSERT INTO Library(Name) VALUES($1)", adminName)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message" : "Admin added successfully"})

		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message" : "password is not correct"})
		}
	}
}


// DELETE Shift Data of any Admin to another Admin
func ShiftAdminToAdmin(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		adminID := c.Query("adminID")
		shiftedAdminID := c.Query("shiftedAdminID")

		var adminExists bool
		var shiftedAdminExists bool
		
		err :=  db.QueryRow("SELECT EXISTS (SELECT 1 FROM Library WHERE ID = $1)", adminID).Scan(&adminExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err =  db.QueryRow("SELECT EXISTS (SELECT 1 FROM Library WHERE ID = $1)", shiftedAdminID).Scan(&shiftedAdminExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		if adminExists {
			if shiftedAdminExists {
				_, err := db.Exec("UPDATE BookInventory SET LibID = $1 WHERE LibID = $2", shiftedAdminID, adminID)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
					return
				}

				_, err = db.Exec("UPDATE Users SET LibID = $1 WHERE LibID = $2", shiftedAdminID, adminID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
					return
				}

				_, err = db.Exec("DELETE FROM Library WHERE ID = $1", adminID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Admin ID %s deleted successfully and All Books and Users are shifted to Admin ID %s", adminID, shiftedAdminID)})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Shifted Admin ID, Try Again"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Message": "Invalid Admin ID, Try Again"})
		}
	}
}