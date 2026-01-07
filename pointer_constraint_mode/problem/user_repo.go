package problem

import (
	"context"

	"gorm.io/gorm"
)

// UserRepositoryForValue 定义了使用值接收者的用户仓库方法
// UserRepositoryForValue defines the methods for user repository using value receiver.
type UserRepositoryForValue interface {
	Create(ctx context.Context, user *UserWithValueMethods) error
	GetByID(ctx context.Context, id uint64) (*UserWithValueMethods, error)
	Update(ctx context.Context, user *UserWithValueMethods) error
	Delete(ctx context.Context, id uint64) error
}

type userRepository struct {
	repo GenericRepositoryForValue[UserWithValueMethods]
}

func NewUserRepository(db *gorm.DB) UserRepositoryForValue {
	return &userRepository{
		//这里我们可以让UserWithValueMethods = T或者*UserWithValueMethods = T,在go语言设计中,当我们使用值方法实现interface时,
		//值类型和指针类型都可以实现该interface
		//Here we can set UserWithValueMethods = T or *UserWithValueMethods = T. In Go's design,
		// when we implement an interface using value methods,
		// both value types and pointer types can implement that interface.
		// 这里推荐使用UserWithValueMethods = T,因为可以避免在使用T时出现空指针异常
		// Using UserWithValueMethods = T is recommended here because it avoids null pointer exceptions when using T.
		repo: NewGenericRepositoryForValue[UserWithValueMethods](db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *UserWithValueMethods) error {
	return r.repo.Create(ctx, user)
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*UserWithValueMethods, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, user *UserWithValueMethods) error {
	return r.repo.Update(ctx, user)
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.repo.Delete(ctx, id)
}

// 在go的设计特性中,当我们使用指针方法实现interface时,实际上只有指针类型可以实现该interface,而值类型不能实现该interface
// In Go's design features, when we implement an interface using pointer methods,
// only pointer types can actually implement that interface; value types cannot.
type UserRepositoryForPointer interface {
	Create(ctx context.Context, user *UserWithPointerMethods) error
	GetByID(ctx context.Context, id uint64) (*UserWithPointerMethods, error)
	Update(ctx context.Context, user *UserWithPointerMethods) error
	Delete(ctx context.Context, id uint64) error
}

type userRepositoryForPointer struct {
	//我们这里只能让*UserWithPointerMethods = T,而不是UserWithPointerMethods = T
	//Here we can only set *UserWithPointerMethods = T, not UserWithPointerMethods = T
	//这很容易导致我们在使用T时出现空指针异常，导致panic
	//For example, if we use UserWithPointerMethods = T,
	//when we call user.GetID(), it will cause a panic because user is a nil pointer.
	repo GenericRepositoryForPointer[*UserWithPointerMethods]
}

func NewUserRepositoryForPointer(db *gorm.DB) UserRepositoryForPointer {
	return &userRepositoryForPointer{
		repo: NewGenericRepositoryForPointer[*UserWithPointerMethods](db),
	}
}

func (r *userRepositoryForPointer) Create(ctx context.Context, user *UserWithPointerMethods) error {
	return r.repo.Create(ctx, user)
}

func (r *userRepositoryForPointer) GetByID(ctx context.Context, id uint64) (*UserWithPointerMethods, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *userRepositoryForPointer) Update(ctx context.Context, user *UserWithPointerMethods) error {
	return r.repo.Update(ctx, user)
}

func (r *userRepositoryForPointer) Delete(ctx context.Context, id uint64) error {
	return r.repo.Delete(ctx, id)
}
