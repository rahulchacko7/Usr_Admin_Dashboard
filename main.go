package main

import (
	db "admin/DB"
	handlers "admin/Handlers"
	"admin/models"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
	}

	router := gin.New()
	db.Db, err = gorm.Open(postgres.Open(os.Getenv("DBS")), &gorm.Config{})
	if err != nil {
		fmt.Println("database not loaded")
		return
	}

	db.Db.AutoMigrate(&models.User{})
	router.LoadHTMLGlob("templates/*.html")
	//	router.Static("/static", "./static")

	//user
	router.GET("/", handlers.LoginHandler)
	router.POST("/", handlers.LoginPost)
	router.GET("/signup", handlers.SignupHandler)
	router.POST("/signup", handlers.SignupPost)
	router.GET("/home", handlers.HomeHandler)
	router.GET("/logout", handlers.LogoutHandler)

	//admin
	router.GET("/admin", handlers.AdminHome)
	router.POST("/admin", handlers.AdminAddUser)
	router.GET("/logoutadmin", handlers.LogoutadminHandler)
	router.GET("/adminupdate", handlers.AdminUpdate)
	router.POST("/adminupdatepost", handlers.AdminUpdatePost)
	router.GET("/admindelete", handlers.AdminDelete)

	// Start the server
	router.Run(":8080")
}
