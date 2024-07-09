#!/bin/bash

echo "Setting up go ------"
go mod vendor
go mod tidy

# Copy sample.env to .env
cp sample.env .env

# MySQL Setup
echo "Please enter your MySQL root password: "
read -s MYSQL_PASSWORD


echo "Creating the database ------------------"
mysql -u root -p$MYSQL_PASSWORD << EOF
CREATE DATABASE IF NOT EXISTS library_management;
USE library_management;
EOF

# Run migrations
echo "Enter your MySQL username:"
read MYSQL_USERNAME

echo "Enter your MySQL password for the user $MYSQL_USERNAME:"
read -s MYSQL_USER_PASSWORD

echo "Enter your MySQL port:"
read -s MYSQL_PORT

# Check if golang-migrate is installed
if ! command -v migrate &> /dev/null
then
    echo "golang-migrate could not be found. Please install it first."
    exit 1
fi

migrate -path database/migration/ -database "mysql://$MYSQL_USERNAME:$MYSQL_USER_PASSWORD@tcp(localhost:$MYSQL_PORT)/library_management" -verbose up


go run ./cmd/main.go