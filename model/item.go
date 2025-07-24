package model

type Item struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"not null"`
	Description string   `gorm:"not null"`
	Price       float64  `gorm:"not null"`
	CategoryID  uint     `gorm:"not null"`
	User        User     `gorm:"foreignKey:UserID;references:ID"`
	Category    Category `gorm:"foreignKey:CategoryID;references:ID"`
	Reviews     []Review `gorm:"foreignKey:ItemID;references:ID"`
	Orders      []Order  `gorm:"foreignKey:ItemID;references:ID"`
}
