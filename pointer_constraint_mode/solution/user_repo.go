package solution

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint64) error
}

type userRepository struct {
	// 这里传递了User和*User,确保*User类型可以实现PointerModel接口
	// Here we pass User and *User, ensuring that only *User types can implement the PointerModel interface.
	// 之后我们可以安全使用Model接口的方法,并且几乎没有额外开销(例如反射和类型断言)
	// We can then safely use the methods of the Model interface with almost no additional overhead
	// (such as reflection and type predicates).
	repo GenericRepository[User, *User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		repo: NewGenericRepository[User](db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.repo.Create(ctx, user)
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*User, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *userRepository) Update(ctx context.Context, user *User) error {
	return r.repo.Update(ctx, user)
}

func (r *userRepository) Delete(ctx context.Context, id uint64) error {
	return r.repo.Delete(ctx, id)
}
