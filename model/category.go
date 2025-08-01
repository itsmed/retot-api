package model

type Category struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"not null;unique"`
	Description string  `gorm:"not null"`
	Items       []Item  `gorm:"foreignKey:CategoryID"`
	Orders      []Order `gorm:"foreignKey:CategoryID"`
}
