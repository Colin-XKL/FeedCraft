package dao

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

// User 程序存储到数据库的用户表
type User struct {
	Username     string `gorm:"primaryKey"`
	NickName     string
	Email        string
	PasswordHash string `gorm:"column:password_hash"` // 这里是密码明文md5 与salt值合并之后再进行sha256计算得到的最终结果
	Salt         string
}

// CreateUser creates a new User record
func CreateUser(db *gorm.DB, user *User, md5Password string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	user.Salt = salt
	user.PasswordHash = HashPasswordWithSalt(md5Password, salt)
	return db.Create(user).Error
}

// GetUserByUsername retrieves a User record by its username
func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUserInfo updates an existing User record
func UpdateUserInfo(db *gorm.DB, user *User) error {
	return db.Model(user).Select("NickName", "Email").Updates(user).Error
}

func UpdateUserPassword(db *gorm.DB, user *User, md5Password string) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	user.Salt = salt
	user.PasswordHash = HashPasswordWithSalt(md5Password, salt)
	return db.Model(user).Select("PasswordHash", "Salt").Updates(user).Error
}

func UpdateUser(db *gorm.DB, user *User, md5Password string) error {
	if md5Password != "" {
		return UpdateUserPassword(db, user, md5Password)
	}
	return UpdateUserInfo(db, user)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return users, nil
}
func generateSalt() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashPasswordWithSalt(password, salt string) string {
	combined := password + salt
	sha256Hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(sha256Hash[:])
}
