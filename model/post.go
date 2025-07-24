package model

type Post struct {
	ID       uint      `gorm:"primaryKey"`
	Title    string    `gorm:"not null"`
	Body     string    `gorm:"not null"`
	UserID   uint      `gorm:"not null"`
	Likes    []Like    `gorm:"foreignKey:PostID;references:ID"`
	Comments []Comment `gorm:"foreignKey:PostID;references:ID"`
}
