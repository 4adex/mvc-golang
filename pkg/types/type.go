package types

import "time"


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

type Transaction struct {
    TransactionID int       `json:"transaction_id"`
    UserID        int       `json:"user_id"`
    BookID        int       `json:"book_id"`
    Status        string    `json:"status"`
    CheckoutTime  time.Time `json:"checkout_time"`
    CheckinTime   time.Time `json:"checkin_time"`
}