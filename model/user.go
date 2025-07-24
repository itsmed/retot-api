package model

// User represents a user in the system
type User struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"unique;not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Posts     []Post    `gorm:"foreignKey:UserID;references:ID"`
    Likes     []Like    `gorm:"foreignKey:UserID;references:ID"`
    Comments  []Comment `gorm:"foreignKey:UserID;references:ID"`
}

