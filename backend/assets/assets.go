package assets


type Library struct {
	ID   int    `json:"ID"`
	Name string `json:"Name"`
}

type User struct {
	ID            int    `json:"ID"`
	Name          string `json:"Name"`
	Email         string `json:"Email"`
	ContactNumber string `json:"ContactNumber"`
	Role          string `json:"Role"`
	LibID         int    `json:"LibID"`
}

type BookInventory struct {
	ISBN           string `json:"ISBN"`
	LibID          int    `json:"LibID"`
	Title          string `json:"Title"`
	Authors        string `json:"Authors"`
	Publisher      string `json:"Publisher"`
	Version        string `json:"Version"`
	TotalCopies    int    `json:"TotalCopies"`
	AvailableCopies int   `json:"AvailableCopies"`
}

type PendingRequest struct {
	ReqID        int       `json:"ReqID"`
	BookID       string    `json:"BookID"`
	ReaderID     int       `json:"ReaderID"`
	RequestDate  string    `json:"RequestDate"`
	RequestType  string    `json:"RequestType"`
}


type RequestEvents struct {
	ReqID        int       `json:"ReqID"`
	BookID       string    `json:"BookID"`
	ReaderID     int       `json:"ReaderID"`
	RequestDate  string    `json:"RequestDate"`
	ApprovalDate string    `json:"ApprovalDate"`
	ApproverID   int       `json:"ApproverID"`
	RequestType  string    `json:"RequestType"`
}

type PendingBooks struct {
	IssueID           int       `json:"IssueID"`
	ISBN              string    `json:"ISBN"`
	ReaderID          int       `json:"ReaderID"`
	IssueApproverID   int       `json:"IssueApproverID"`
	IssueStatus       string    `json:"IssueStatus"`
	IssueDate         string    `json:"IssueDate"`
	ExpectedReturnDate string   `json:"ExpectedReturnDate"`
}

type IssueRegistry struct {
	IssueID           int       `json:"IssueID"`
	ISBN              string    `json:"ISBN"`
	ReaderID          int       `json:"ReaderID"`
	IssueApproverID   int       `json:"IssueApproverID"`
	IssueStatus       string    `json:"IssueStatus"`
	IssueDate         string    `json:"IssueDate"`
	ExpectedReturnDate string   `json:"ExpectedReturnDate"`
	ReturnDate        string    `json:"ReturnDate"`
	ReturnApproverID  int       `json:"ReturnApproverID"`
}

func DateFunc(t string) string {
	return t[0:10]
}