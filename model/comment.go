package model

type Comment struct {
	ID     uint   `gorm:"primaryKey"`
	Body   string `gorm:"not null"`
	UserID uint   `gorm:"not null"`
	PostID uint   `gorm:"not null"`
	User   User   `gorm:"foreignKey:UserID;references:ID"`
}
