package main

import (
	"github.com/sudaratimjai/finalexam/customer"
	"github.com/sudaratimjai/finalexam/database"
)

func main() {
	database.Conn()
	customer.CreateTable()
	r := customer.NewRouter()
	r.Run(":2019")
}
