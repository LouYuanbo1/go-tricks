package model

type User struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
}

func (u *User) GetID() uint64 {
	return u.ID
}

func (u *User) GetPrimaryKey() string {
	return "id"
}

func (u *User) TableName() string {
	return "users"
}
