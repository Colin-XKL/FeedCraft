package dao

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"gorm.io/gorm"
)

type User struct {
	Username     string `gorm:"primaryKey"`
	NickName     string
	Email        string
	PasswordHash string `gorm:"column:password_hash"`
	Salt         string
}

type UserInfo struct {
	Username string `json:"username"`
	NickName string `json:"nickname"`
	Email    string `json:"email"`
}

// CreateUser creates a new User record
func CreateUser(db *gorm.DB, user *User) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	user.Salt = salt
	user.PasswordHash = hashPassword(user.Password, salt)
	user.Password = "" // Clear the password field
	return db.Omit("Password").Create(user).Error
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
		salt, err := generateSalt()
		if err != nil {
			return err
		}
		user.Salt = salt
		user.PasswordHash = hashPassword(user.Password, salt)
		user.Password = "" // 清空密码
	}
	return db.Omit("Password").Save(user).Error
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
func generateSalt() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func hashPassword(password, salt string) string {
	md5Hash := md5.Sum([]byte(password))
	md5Password := hex.EncodeToString(md5Hash[:])
	combined := md5Password + salt
	sha256Hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(sha256Hash[:])
}
