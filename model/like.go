package model

type Like struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"not null"`
	PostID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID;references:ID"`
	Post   Post `gorm:"foreignKey:PostID;references:ID"`
}
