package problem

type Model interface {
	GetID() uint64
	GetPrimaryKey() string
	TableName() string
}

// UserWithPointerMethods 实现了 Model 接口，但是 GetID 方法是指针接收者
// UserWithPointerMethods implements the Model interface, but the GetID method is a pointer receiver.
type UserWithPointerMethods struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
}

func (u *UserWithPointerMethods) GetID() uint64 {
	return u.ID
}

func (u *UserWithPointerMethods) GetPrimaryKey() string {
	return "id"
}

func (u *UserWithPointerMethods) TableName() string {
	return "users"
}

// UserWithValueMethods 实现了 Model 接口，但是 GetID 方法是值接收者
// UserWithValueMethods implements the Model interface, but the GetID method is a value receiver.
type UserWithValueMethods struct {
	ID    uint64 `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Email string `gorm:"not null;unique"`
}

func (u UserWithValueMethods) GetID() uint64 {
	return u.ID
}

func (u UserWithValueMethods) GetPrimaryKey() string {
	return "id"
}

func (u UserWithValueMethods) TableName() string {
	return "users"
}
