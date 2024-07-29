package dao

import (
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"primaryKey"`
	NickName string
	Password string
	Email    string
}

// CreateUser creates a new User record
func CreateUser(db *gorm.DB, user *User) error {
	return db.Create(user).Error
}

// GetUserByUsername retrieves a User record by its username
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	return &user, result.Error
}

// UpdateUser updates an existing User record
func UpdateUser(db *gorm.DB, user *User) error {
	return db.Save(user).Error
}

// DeleteUser deletes a User record by its username
func DeleteUser(db *gorm.DB, username string) error {
	var user User
	return db.Where("username = ?", username).Delete(&user).Error
}

// ListUsers retrieves all User records
func ListUsers(db *gorm.DB) ([]*User, error) {
	var users []*User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
