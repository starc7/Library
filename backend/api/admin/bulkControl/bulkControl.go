package bulkControl

import (
	"encoding/csv"
	"fmt"
	"io"
	"database/sql"
	"os"
	
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
)

func ImportBooksInBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
			return
		}
		defer file.Close()

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
			return
		}
		defer tx.Rollback()

		reader := csv.NewReader(file)
		var values []string

		_, err = reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read CSV data"})
				return
			}

			// Assuming CSV columns are in the order: ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies
			isbn := record[0]
			libID, _ := strconv.Atoi(record[1])
			title := record[2]
			authors := record[3]
			publisher := record[4]
			version := record[5]
			totalCopies, _ := strconv.Atoi(record[6])
			availableCopies, _ := strconv.Atoi(record[7])

			// Escape any single quotes in the data
			title = strings.Replace(title, "'", "''", -1)
			authors = strings.Replace(authors, "'", "''", -1)
			publisher = strings.Replace(publisher, "'", "''", -1)
			version = strings.Replace(version, "'", "''", -1)

			// Construct the SQL statement for bulk insert
			values = append(values, fmt.Sprintf("('%s', %d, '%s', '%s', '%s', '%s', %d, %d)", isbn, libID, title, authors, publisher, version, totalCopies, availableCopies))

			var isbnExists bool
			err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM BOOKINVENTORY WHERE ISBN = $1)",isbn).Scan(&isbnExists)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}


			if isbnExists {
				_, err = db.Exec("UPDATE BookInventory SET AvailableCopies = AvailableCopies + $1, TotalCopies = TotalCopies + $2 WHERE ISBN = $3", availableCopies, totalCopies, isbn)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
					return
				}
				values = nil
			} else {
				stmt := fmt.Sprintf("INSERT INTO BookInventory (ISBN, LibID, Title, Authors, Publisher, Version, TotalCopies, AvailableCopies) VALUES %s", strings.Join(values, ","))
				_, err = tx.Exec(stmt)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into the database"})
					return
				}

				values = nil
			}
		}

		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "CSV data imported successfully"})
	}	
}


func ExportBooksInBulk(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tableName := c.Param("tableName")
		adminID := c.Param("adminID")

		filePath := fmt.Sprintf("%s.csv", tableName)
		file, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create CSV file"})
			return
		}
		defer file.Close()

		rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE LibID = %s", tableName, adminID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query data from table"})
			return
		}
		defer rows.Close()

		writer := csv.NewWriter(file)

		// Write the column headers to the CSV file
		columns, err := rows.Columns()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get column headers"})
			return
		}
		writer.Write(columns)

		// Fetch rows and write to CSV file
		for rows.Next() {
			// Retrieve row data
			values := make([]interface{}, len(columns))
			for i := range values {
				values[i] = new(interface{})
			}

			err := rows.Scan(values...)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row data"})
				return
			}

			// Convert values to strings and write to CSV
			var row []string
			for _, v := range values {
				row = append(row, fmt.Sprintf("%v", *v.(*interface{})))
			}
			writer.Write(row)
		}

		// Flush the writer to write the remaining data to file
		writer.Flush()

		// Check for any writer error
		if err := writer.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data to CSV file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Data exported to CSV successfully"})
	}
}