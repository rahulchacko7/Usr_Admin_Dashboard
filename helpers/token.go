package helpers

import (
	"admin/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(user models.User, c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": user.Role,
		"user": user.UserName,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("KEY")))

	if err == nil {
		fmt.Println("token created")
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("cookie", tokenString, 3600, "", "", false, true)
}
