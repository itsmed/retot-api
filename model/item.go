package model

type Item struct {
	ID          uint     `gorm:"primaryKey"`
	Name        string   `gorm:"not null"`
	Description string   `gorm:"not null"`
	Price       float64  `gorm:"not null"`
	UserID      uint     `gorm:"not null"`
	User        User     `gorm:"foreignKey:UserID"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Reviews     []Review `gorm:"foreignKey:ItemID"`
	Orders      []Order  `gorm:"foreignKey:ItemID"`
}
