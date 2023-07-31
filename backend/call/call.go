package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"champ.com/connectDB"
	"champ.com/api/owner"
	"champ.com/api/admin/bookControl"
	"champ.com/api/admin/requestControl"
	"champ.com/api/admin/userControl"
	"champ.com/api/admin/issueControl"
	"champ.com/api/user"
	"champ.com/api/admin/bulkControl"
)

func main() {
	db, err := connectDB.GetDB()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("api/owner/getAllAdminData", owner.GetAllAdminData(db))
	router.POST("api/owner/addNewAdmin", owner.AddNewAdmin(db))
	router.GET("api/owner/getAllBooks", owner.GetAllBookData(db))
	router.POST("api/admin/bookControl/addNewBook", bookControl.AddNewBookToLib(db))
	router.GET("api/admin/bookControl/getAllBook", bookControl.GetAllBooksByAdmin(db))
	router.PUT("api/admin/bookControl/updateBook", bookControl.UpdateBookRecord(db))
	router.POST("api/admin/userControl/addNewUser", userControl.AddNewUser(db))
	router.GET("api/admin/userControl/getAllUsers", userControl.GetAllUserByAdmin(db))
	router.GET("api/user/searchBook/byAuthors", user.SearchBookByAuthorsName(db))
	router.GET("api/user/searchBook/byTitle", user.SearchBookByTitle(db))
	router.GET("api/user/searchBook/byPublisher", user.SearchBookByPublisher(db))
	router.POST("api/user/raiseRequest", user.RaiseRequestForBook(db))
	router.GET("api/admin/requestControl/getAllPendingRequest", requestControl.GetAllPendingRequest(db))
	router.DELETE("api/admin/requestControl/rejectRequest", requestControl.RejectRequest(db))
	router.PUT("api/admin/requestControl/approveRequest", requestControl.ApproveRequstByAdmin(db))
	router.GET("api/user/getIssueRegistry", user.GetIssueRegistryByUser(db))
	router.GET("api/admin/issueControl/getPendingBooks", issueControl.GetAllPendingBooksDetails(db))
	router.PUT("api/admin/bookControl/returnBook", bookControl.BookReturnByUser(db))
	router.GET("api/admin/issueControl/getAllIssues", issueControl.GetAllIssueRegistry(db))
	router.GET("api/admin/requestControl/getAllRequest", requestControl.GetAllRequests(db))
	router.DELETE("api/owner/removeAdmin", owner.ShiftAdminToAdmin(db))


	router.POST("api/importBooksInBulk", bulkControl.ImportBooksInBulk(db))
	router.GET("/exportCSV/:tableName/admin/:adminID", bulkControl.ExportBooksInBulk(db))

	router.Run("localhost:8080")
}