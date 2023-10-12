package db

import (
	"admin/models"

	"gorm.io/gorm"
)

var Db *gorm.DB
var UserList []models.User
