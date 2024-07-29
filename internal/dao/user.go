package dao

type User struct {
	Username string `gorm:"primaryKey"`
	NickName string
	Password string
	Email    string
}
