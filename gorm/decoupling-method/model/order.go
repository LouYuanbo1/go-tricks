package model

type Order struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null"`
	ProductID uint64 `gorm:"not null"`
}

func (o *Order) GetID() uint64 {
	return o.ID
}

func (o *Order) GetPrimaryKey() string {
	return "id"
}

func (o *Order) TableName() string {
	return "orders"
}
