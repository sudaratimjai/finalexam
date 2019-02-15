package customer

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sudaratimjai/finalexam/database"
)

var customers []Customer

func CreateTable() {

	ctb := `
	CREATE TABLE IF NOT EXISTS customers(
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT,
		status TEXT
	);
	`
	_, err := database.Conn().Exec(ctb)

	if err != nil {
		log.Fatal("can't create customer table", err)
		return
	}

	fmt.Println("create customer table sucess")
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(loginMiddleware)
	r.POST("/customers", insertCustomerHandler)
	r.GET("/customers/:id", getCustomerByIDHandler)
	r.GET("/customers", getAllCustomerHandler)
	r.PUT("/customers/:id", updateCustomerHandler)
	r.DELETE("/customers/:id", deleteCustomerHandler)

	return r
}

func insertCustomerHandler(c *gin.Context) {

	var item Customer
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	row := database.InsertCustomer(item.Name, item.Email, item.Status)

	cus := Customer{}
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		log.Fatal("can't Scan row customer :", err)
		return
	}

	c.JSON(http.StatusCreated, cus)
}

func getCustomerByIDHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := database.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "prepare select error"+err.Error())
		return
	}

	row := stmt.QueryRow(id)
	fmt.Println("id", id)
	cus := Customer{}
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		log.Fatal("can't Scan row into vaiable", err)
		return
	}

	c.JSON(http.StatusOK, cus)
}

func getAllCustomerHandler(c *gin.Context) {

	stmt, err := database.GetAllCustomer()
	if err != nil {
		log.Fatal("can't prepare query all customers", err)
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("can't query all customers", err)
		return
	}

	var customers = []Customer{}

	for rows.Next() {
		c := Customer{}
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Status)
		if err != nil {
			log.Fatal("can't scan row into variable", err)
			return
		}

		customers = append(customers, c)
	}
	c.JSON(http.StatusOK, customers)
}

func updateCustomerHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var item Customer
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	stmt, err := database.UpdateCustomer()
	if err != nil {
		log.Fatal("can't update customer", err)
		return
	}

	if _, err := stmt.Exec(id, item.Name, item.Email, item.Status); err != nil {
		log.Fatal("error execute update", err)
		return
	}

	stmt, err = database.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "prepare error"+err.Error())
		return
	}

	row := stmt.QueryRow(id)

	cus := Customer{}
	err = row.Scan(&cus.ID, &cus.Name, &cus.Email, &cus.Status)
	if err != nil {
		log.Fatal("can't Scan row into vaiable", err)
		return
	}

	c.JSON(http.StatusOK, cus)
}

func deleteCustomerHandler(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	stmt, err := database.DeleteCustomer()
	if err != nil {
		log.Fatal("can't delete customer", err)
		return
	}

	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("error execute delete", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func loginMiddleware(c *gin.Context) {

	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Status code is 401 Unauthorized")
		c.Abort()
		return
	}
	c.Next()
	log.Panicln("ending")
}
