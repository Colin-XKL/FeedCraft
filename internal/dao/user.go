package dao

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
func CreateUser(db *gorm.DB, user *User) error {
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	user.Salt = salt
	user.PasswordHash = HashPassword(user.Password, salt)
	user.Password = "" // 清空密码
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

// UpdateUser updates an existing User record
func UpdateUser(db *gorm.DB, user *User) error {
	if user.Password != "" {
		salt, err := generateSalt()
		if err != nil {
			return err
		}
		user.Salt = salt
		user.PasswordHash = HashPassword(user.Password, salt)
		user.Password = "" // 清空密码
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

func HashPassword(password, salt string) string {
	combined := password + salt
	sha256Hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(sha256Hash[:])
}
