package solution

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type GenericRepository[T any, PT PointerModel[T]] interface {
	Create(ctx context.Context, model *T) error
	GetByID(ctx context.Context, id uint64) (*T, error)
	Update(ctx context.Context, model *T) error
	Delete(ctx context.Context, id uint64) error
}

type genericRepository[T any, PT PointerModel[T]] struct {
	db *gorm.DB
}

func NewGenericRepository[T any, PT PointerModel[T]](db *gorm.DB) GenericRepository[T, PT] {
	return &genericRepository[T, PT]{db: db}
}

func (r *genericRepository[T, PT]) Create(ctx context.Context, model *T) error {
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepository[T, PT]) GetByID(ctx context.Context, id uint64) (*T, error) {
	var model T
	//&model：首先获取值类型变量 model 的内存地址，得到一个 *T 类型的指针
	//PT(&model)：将 *T 类型的指针转换为 PT 类型（PointerModel[T]）
	//&model: First, retrieve the memory address of the value type variable model, resulting in a pointer of type *T.
	//PT(&model): Convert the pointer of type *T to type PT (PointerModel[T]).
	//确保在使用指针方法时,不会出现空指针异常
	//Ensure that when using pointer methods, there will be no nil pointer exceptions.
	pt := PT(&model)
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", pt.GetPrimaryKey()), id).
		First(&model).
		Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *genericRepository[T, PT]) Update(ctx context.Context, model *T) error {
	err := r.db.WithContext(ctx).Updates(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *genericRepository[T, PT]) Delete(ctx context.Context, id uint64) error {
	var model T
	pt := PT(&model)
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s = ?", pt.GetPrimaryKey()), id).
		Delete(&model).
		Error
	if err != nil {
		return err
	}
	return nil
}
