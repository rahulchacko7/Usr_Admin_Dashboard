package handlers

import (
	db "admin/DB"
	"admin/helpers"
	"fmt"
	"regexp"

	"admin/midleware"
	"admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	// validate cookie
	ok := midleware.ValidateCookies(c)
	if !ok {
		c.HTML(http.StatusOK, "signup.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/")
}

func SignupPost(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	var errors models.Invalid
	userName := c.Request.FormValue("Name")
	userEmail := c.Request.FormValue("Email")

	if userName == "" {
		errors.NameError = "name should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	if userEmail == "" {
		errors.EmailError = "Email should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	if !regex.MatchString(userEmail) {
		errors.EmailError = "Email not in the correct format"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	password := c.Request.FormValue("Password")
	confirmPassword := c.Request.FormValue("ConfirmPassword")
	if password == "" {
		errors.PasswordError = "Password should not be empty"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}
	if password != confirmPassword {
		errors.PasswordError = "password does not match"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	//check if user already exists
	var count int
	if err := db.Db.Raw("SELECT COUNT(*) FROM users WHERE email=$1", userEmail).Scan(&count).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "signup.html", nil)
		return
	}
	if count > 0 {
		errors.EmailError = "user already exists"
		c.HTML(http.StatusBadRequest, "signup.html", errors)
		return
	}

	//insert data into database
	err := db.Db.Exec("INSERT INTO users(user_name,email,password) VALUES($1,$2,$3)", userName, userEmail, password).Error
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusOK, "signup.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/")

	// user := models.User{
	// 	Role:     "user",
	// 	UserName: userName,
	// }
	// helpers.CreateToken(user, c)
	// c.Redirect(http.StatusFound, "/home")
	// return

}

func LoginHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	// validate cookies
	ok := midleware.ValidateCookies(c)
	role, _, _ := midleware.FindRole(c)
	if !ok {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		if role == "user" {
			c.Redirect(http.StatusFound, "/home")
			return
		} else if role == "admin" {
			c.Redirect(http.StatusFound, "/admin")
			return
		}
		c.HTML(http.StatusBadRequest, "login.html", nil)
	}
}

func LoginPost(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	var error models.Invalid
	Newmail := c.Request.FormValue("Email")
	Newpassword := c.Request.FormValue("Password")
	var compare models.Compare

	if err := db.Db.Raw("SELECT password,role,user_name FROM users WHERE email=$1", Newmail).Scan(&compare).Error; err != nil {
		fmt.Println(err)
		c.HTML(http.StatusBadRequest, "login.html", nil)
		return
	}

	if compare.Password != Newpassword {
		error.PasswordError = "check password again"
		c.HTML(http.StatusBadRequest, "login.html", error)
		return
	}
	if compare.Role == "user" {
		user := models.User{
			Role:     compare.Role,
			UserName: compare.UserName,
		}
		helpers.CreateToken(user, c)
		c.Redirect(http.StatusFound, "/home")
		return
	} else if compare.Role == "admin" {
		user := models.User{
			Role:     compare.Role,
			UserName: compare.UserName,
		}
		helpers.CreateToken(user, c)
		c.Redirect(http.StatusFound, "/admin")
		return
	} else {
		error.EmailError = "role mismatch"
		c.HTML(http.StatusOK, "login.html", error)
		return
	}
}

// var role,User string
func HomeHandler(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Expires", "0")
	ok := midleware.ValidateCookies(c)
	role, User, _ := midleware.FindRole(c)
	if !ok {
		c.HTML(http.StatusOK, "login.html", nil)
	} else {
		if role == "user" {
			c.HTML(http.StatusOK, "home.user.html", gin.H{"UserName": User})
			return
		} else {
			c.Redirect(http.StatusFound, "/")
			return
		}

	}
}

func LogoutHandler(c *gin.Context) {
	midleware.DeleteCookie(c)
	c.Redirect(http.StatusFound, "/")
}
