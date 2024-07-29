package dao

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `gorm:"primaryKey"`
	NickName string
	Email    string
	PasswordHash []byte `gorm:"column:password_hash"`
}

type UserInfo struct {
	Username string `json:"username"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

// CreateUser creates a new User record
func CreateUser(db *gorm.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword
	return db.Create(user).Error
}

// GetUserByUsername retrieves a User record by its username
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser updates an existing User record
func UpdateUser(db *gorm.DB, user *User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = hashedPassword
	}
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
