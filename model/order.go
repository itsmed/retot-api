package model

// Order represents an order for an item
type Order struct {
	ID         uint    `gorm:"primaryKey"`
	ItemID     uint    `gorm:"not null"`
	UserID     uint    `gorm:"not null"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"not null"`
	CategoryID uint    `gorm:"not null"`
	Item       Item    `gorm:"foreignKey:ItemID;references:ID"`
	User       User    `gorm:"foreignKey:UserID;references:ID"`
}
