package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type person struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	City      string `json:"city"`
}

func main() {
	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function
	//db, err = gorm.Open("sqlite3", "./gorm.db")
	db, _ = gorm.Open("mysql", "root:db@tcp(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	db.AutoMigrate(&person{})
	r := gin.Default()
	r.GET("/people/", getPeople)
	r.GET("/people/:id", getPerson)
	r.POST("/people", createPerson)
	r.PUT("/people/:id", updatePerson)
	r.DELETE("/people/:id", deletePerson)
	r.Run(":8080")
}

func deletePerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person person
	d := db.Where("id = ?", id).Delete(&person)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}

func updatePerson(c *gin.Context) {
	var person person
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	c.BindJSON(&person)
	db.Save(&person)
	c.JSON(200, person)
}

func createPerson(c *gin.Context) {
	var person person
	c.BindJSON(&person)
	db.Create(&person)
	c.JSON(200, person)
}

func getPerson(c *gin.Context) {
	id := c.Params.ByName("id")
	var person person
	if err := db.Where("id = ?", id).First(&person).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, person)
	}
}

func getPeople(c *gin.Context) {
	var people []person
	if err := db.Find(&people).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, people)
	}
}
