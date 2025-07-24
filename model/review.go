package model

// Review represents a review for an item
type Review struct {
	ID        uint   `gorm:"primaryKey"`
	ItemID    uint   `gorm:"not null"`
	UserID    uint   `gorm:"not null"`
	Rating    int    `gorm:"not null"` // Rating out of 5
	Comment   string `gorm:"not null"`
	Item      Item   `gorm:"foreignKey:ItemID;references:ID"`
	User      User   `gorm:"foreignKey:UserID;references:ID"`
}
