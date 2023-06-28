package models

import "time"

type Calculate struct {
	OriginID        int    `gorm:"originid" binding:"required" json:"originid"`
	DestinationID   int    `gorm:"destinationid" json:"destinationid"`
	LatOrigin       string `gorm:"lat_origin" json:"lat_origin"`
	LongOrigin      string `gorm:"long_origin" json:"long_origin"`
	LatDestination  string `gorm:"lat_destination" json:"lat_destination"`
	LongDestination string `gorm:"long_destination" json:"long_destination"`
	//
	OtherStatus  int       `gorm:"other_status" json:"other_status"`
	OtherID      int       `gorm:"otherid" json:"otherid"`
	LatOther     string    `gorm:"lat_other" json:"lat_other"`
	LongLatOther string    `gorm:"long_other" json:"long_other"`
	CreatedAt    time.Time `gorm:"created_at" `
	UpdatedAt    time.Time `gorm:"updated_at"`
}

func (Calculate) TableName() string {
	return "calculates"
}
