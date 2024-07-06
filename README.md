# mvc-golang

- Clone the repo. From the root directory of the cloned repo:
```
- go mod vendor
- go mod tidy
- cp sample.env .env
```

- MYSQL Setup:
1. `mysql -u root -p` : and enter password
2. Create a new database 'books': `CREATE DATABASE library_management;`
3. Connect to the database: `USE library_management;`

- Do all the migrations:
1. Ensure that you have golang migrate installed.
2. Change the username and password and run `migrate -path database/migration/ -database "mysql://username:password@tcp(localhost:3306)/library_management" -verbose up` 

- Running the server:
1. `go build -o mvc ./cmd/main.go`
2.  Run the binary file: `./mvc`