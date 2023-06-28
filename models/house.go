package models

import "time"

type House struct {
	ID        uint      `gorm:"id"`
	Tipe      string    `gorm:"tipe" binding:"required"`
	Alamat    string    `gorm:"alamat" binding:"required"`
	Lat       string    `gorm:"lat" binding:"required"`
	Long      string    `gorm:"long" binding:"required"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
	DeletedAt time.Time `gorm:"deleted_at"`
}

func (House) TableName() string {
	return "houses"
}
