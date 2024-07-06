package types

import (
	"database/sql"
	// "time"
)


type User struct {
    ID int `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`
    Salt string `json:"salt"`
    RequestStatus string `json:"request_status"`
}

type Book struct {
    ID             int    `json:"id"`
    Title          string `json:"title"`
    Author         string `json:"author"`
    ISBN           string `json:"isbn"`
    PublicationYear int    `json:"publication_year"`
    TotalCopies    uint   `json:"total_copies"`
    AvailableCopies uint  `json:"available_copies"`
}

type Response struct {
    Message string `json:"message"`
    Type string `json:"type"`
}

type AdminRequest struct {
    ID            int    `json:"id"`
    Username      string `json:"username"`
    RequestStatus string `json:"request_status"`
}

type Transaction struct {
    TransactionID string       `json:"transaction_id"`
    UserID        string       `json:"user_id"`
    BookID        string       `json:"book_id"`
    Status        string    `json:"status"`
    CheckoutTime  string `json:"checkout_time"`
    CheckinTime   sql.NullString `json:"checkin_time"`
}

type PendingTransaction struct {
    TransactionID string
    Title         string
    Status        string
    CheckoutTime  string
    CheckinTime   sql.NullString
}

type History struct {
	TransactionID string
	Title         string
	Status        string
	CheckoutTime  string
	CheckinTime   string
}

type Holding struct {
    TransactionID string
    Title         string
    Author        string
    CheckoutTime  string
}