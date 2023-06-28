package models

type User struct {
	ID       int    `gorm:"id"`
	Email    string `gorm:"email"`
	Username string `gorm:"username" `
	Password string `gorm:"password"`
	Token    string `gorm:"token" `
}

func (User) TableName() string {
	return "user"
}
