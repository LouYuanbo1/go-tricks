package problem

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// GenericRepositoryForValue 是一个泛型仓库接口，用于操作值接收者的模型
// GenericRepositoryForValue is a generic repository interface for value receivers.
type GenericRepositoryForValue[T Model] interface {
	Create(ctx context.Context, model *T) error
	GetByID(ctx context.Context, id uint64) (*T, error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, id uint64) error
}

type genericRepositoryForValue[T Model] struct {
	db *gorm.DB
}

func NewGenericRepositoryForValue[T Model](db *gorm.DB) GenericRepositoryForValue[T] {
	return &genericRepositoryForValue[T]{db: db}
}

func (r *genericRepositoryForValue[T]) Create(ctx context.Context, model *T) error {
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepositoryForValue[T]) GetByID(ctx context.Context, id uint64) (*T, error) {
	var model T
	// 注意：这里使用了值接收者的 GetPrimaryKey 方法
	// Note: We use the value receiver's GetPrimaryKey method here.
	// 我们这里的T是一个值类型，所以我们可以直接使用GetPrimaryKey 方法,不用担心空指针问题
	// T is a value type, so we can use GetPrimaryKey method directly, without worrying about nil pointer.
	// 但是如果T对应的结构体非常大的情况下,使用值方法会导致性能问题,尤其是获取结构体内部字段时(例如:model.GetID()),
	// 因为值方法会复制结构体,而指针方法不会
	// However, if the structure corresponding to T is very large, using the value method can lead to performance issues,
	// especially when accessing struct fields(e.g:model.GetID()).Because value method will copy struct.
	// 所以我们一般希望使用指针方法,详见下方注释
	// We generally expect to use pointer methods, see the comments below for details.
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", model.GetPrimaryKey()), id).
		First(&model).
		Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *genericRepositoryForValue[T]) Update(ctx context.Context, model *T) error {
	err := r.db.WithContext(ctx).Updates(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepositoryForValue[T]) Delete(ctx context.Context, id uint64) error {
	var model T
	//这里与上方问题相同
	//This is the same as the question above.
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", model.GetPrimaryKey()), id).
		Delete(&model).
		Error
	if err != nil {
		return err
	}
	return nil
}

// GenericRepositoryForPointer 是一个泛型仓库接口，用于操作指针接收者的模型
// GenericRepositoryForPointer is a generic repository interface for pointer receivers.
type GenericRepositoryForPointer[T Model] interface {
	Create(ctx context.Context, model T) error
	GetByID(ctx context.Context, id uint64) (T, error)
	Update(ctx context.Context, model T) error
	Delete(ctx context.Context, id uint64) error
}

type genericRepositoryForPointer[T Model] struct {
	db *gorm.DB
}

func NewGenericRepositoryForPointer[T Model](db *gorm.DB) GenericRepositoryForPointer[T] {
	return &genericRepositoryForPointer[T]{db: db}
}

func (r *genericRepositoryForPointer[T]) Create(ctx context.Context, model T) error {
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepositoryForPointer[T]) GetByID(ctx context.Context, id uint64) (T, error) {
	var model T
	// 注意：这里使用了指针接收者的 GetPrimaryKey 方法
	// Note: We use the pointer receiver's GetPrimaryKey method here.
	// 我们这里的T是一个指针类型，如果我们直接使用GetPrimaryKey 方法,可能会导致空指针问题,因为这里的model等于nil
	// Here, T is a pointer type. If we directly use the GetPrimaryKey method,
	// it may cause a null pointer exception because model here is equal to nil.
	// 这里暂时不会造成问题,因为go允许调用空指针的方法,如果方法中没有涉及结构体内部字段
	// It will not cause a problem here because Go allows calling methods on nil pointers,
	// if the method does not involve struct fields.
	// 但是如果方法中涉及了结构体内部字段,就会导致空指针问题(例如:model.GetID())
	// However, if the GetPrimaryKey method involves struct fields, it may cause a null pointer exception,
	// e.g: model.GetID()
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", model.GetPrimaryKey()), id).
		First(&model).
		Error
	if err != nil {
		return model, err
	}
	return model, nil
}

func (r *genericRepositoryForPointer[T]) Update(ctx context.Context, model T) error {
	err := r.db.WithContext(ctx).Updates(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepositoryForPointer[T]) Delete(ctx context.Context, id uint64) error {
	var model T
	//这里与上方问题相同
	//This is the same as the question above.
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", model.GetPrimaryKey()), id).
		Delete(&model).
		Error
	if err != nil {
		return err
	}
	return nil
}
