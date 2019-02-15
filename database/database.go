package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Conn() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	return db
}

func InsertCustomer(name, email, status string) *sql.Row {
	return Conn().QueryRow("INSERT INTO customers (name,email,status) values ($1,$2,$3) RETURNING id,name,email,status", name, email, status)
}

func GetCustomerByID(id int) (*sql.Stmt, error) {
	return Conn().Prepare("SELECT id, name, email, status FROM customers where id=$1")
}

func GetAllCustomer() (*sql.Stmt, error) {
	return Conn().Prepare("SELECT id ,name, email, status FROM customers")
}

func UpdateCustomer() (*sql.Stmt, error) {
	return Conn().Prepare("UPDATE customers SET name=$2, email=$3, status=$4 WHERE id=$1;")
}

func DeleteCustomer() (*sql.Stmt, error) {
	return Conn().Prepare("DELETE FROM customers WHERE id=$1")
}
