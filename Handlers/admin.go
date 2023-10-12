package handlers

import (
	db "admin/DB"
	"admin/midleware"
	"admin/models"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AdminResponse struct {
	Name    string
	Users   []models.UserDetails
	Invalid models.Invalid
}

var errors models.Invalid

func AdminHome(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")

	ok := midleware.ValidateCookies(c)
	if !ok {
		c.HTML(http.StatusOK, "login.html", nil)
	}

	role, name, err := midleware.FindRole(c)
	if err != nil {
		fmt.Println(err)
	}
	if role != "admin" {
		c.Redirect(http.StatusFound, "/")
		return
	}

	var Collect []models.UserDetails
	if err := db.Db.Raw("select user_name,email from users").Scan(&Collect).Error; err != nil {
		fmt.Println("could not fetch user details")
	}
	result := AdminResponse{
		Name:    name,
		Users:   Collect,
		Invalid: errors,
	}
	c.HTML(http.StatusOK, "adminhomepage.html", gin.H{
		"title": result,
	})
}

func AdminAddUser(c *gin.Context) {
	fmt.Println("sucess")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	ok := midleware.ValidateCookies(c)
	role, _, _ := midleware.FindRole(c)
	if !ok || role != "admin" {
		c.HTML(http.StatusOK, "login.html", nil)
	}
	fmt.Println("sucess")

	userName := c.Request.FormValue("Name")
	userEmail := c.Request.FormValue("Email")
	userPassword := c.Request.FormValue("Password")

	errors.NameError = ""
	errors.EmailError = ""
	errors.PasswordError = ""
	errors.CommonError = ""

	if userName == "" {
		errors.EmailError = ""
		errors.PasswordError = ""
		errors.CommonError = ""
		errors.NameError = "name should not be empty"
		c.Redirect(http.StatusFound, "/admin")
		return
	} else if userEmail == "" {
		errors.NameError = ""
		errors.PasswordError = ""
		errors.CommonError = ""
		errors.EmailError = "Email should not be empty"
		c.Redirect(http.StatusFound, "/admin")
		return
	} else if userPassword == "" {
		errors.NameError = ""
		errors.EmailError = ""
		errors.CommonError = ""
		errors.PasswordError = "Password should not be empty"
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(userEmail) {
		errors.EmailError = "Email not in the correct format"
		c.Redirect(http.StatusFound, "/admin")
		return
	}
	var count int
	if err := db.Db.Raw("SELECT COUNT(*) FROM users WHERE email=$1", userEmail).Scan(&count).Error; err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	if count > 0 {
		errors.CommonError = "user already exists"
		c.Redirect(http.StatusFound, "/admin")
		return
	}

	var userRole string
	if c.Request.FormValue("checkbox") == "on" {
		userRole = "admin"
	} else {
		userRole = "user"
	}

	err := db.Db.Exec("INSERT INTO users (user_name,role,email,password) VALUES($1,$2,$3,$4)", userName, userRole, userEmail, userPassword).Error
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/admin")
		return
	}
	c.Redirect(http.StatusFound, "/admin")

}

func AdminUpdate(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	ok := midleware.ValidateCookies(c)
	if !ok {
		c.HTML(http.StatusOK, "login.html", nil)
	}
	username := c.Query("Username")
	email := c.Query("Email")
	c.HTML(http.StatusOK, "updateuser.html", gin.H{
		"UserName": username,
		"Email":    email,
	})

}

func AdminUpdatePost(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	fmt.Println("reahed heare post")
	ok := midleware.ValidateCookies(c)
	if !ok {
		c.HTML(http.StatusOK, "login.html", nil)
	}

	email := c.Query("Email")
	userName := c.Request.FormValue("Name")
	err := db.Db.Exec("UPDATE users SET user_name=$1 where email=$2", userName, email).Error
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusFound, "/admin")

}

func AdminDelete(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	ok := midleware.ValidateCookies(c)
	role, _, _ := midleware.FindRole(c)
	if !ok || role != "admin" {
		c.HTML(http.StatusOK, "login.html", nil)
	}
	email := c.Query("Email")
	if err := db.Db.Exec("delete from users where email=?", email).Error; err != nil {
		fmt.Println("could not fetch user details")
	}
	c.Redirect(http.StatusFound, "/admin")
}

func LogoutadminHandler(c *gin.Context) {
	midleware.DeleteCookie(c)
	c.Redirect(http.StatusFound, "/")
}
