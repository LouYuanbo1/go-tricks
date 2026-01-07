package repo

import (
	"context"
	"go-tricks/gorm/decoupling-method/model"

	genericRepo "go-tricks/gorm/decoupling-method/repo/generic"

	"gorm.io/gorm"
)

type repoFactory struct {
	db *gorm.DB
}

func NewRepoFactory(db *gorm.DB) *repoFactory {
	return &repoFactory{db: db}
}

func (f *repoFactory) withTransaction(tx *gorm.DB) *repoFactory {
	return &repoFactory{
		db: tx,
	}
}

// Transaction
func (f *repoFactory) Transaction(ctx context.Context, fn func(factory RepoFactory) error) error {
	// 使用gorm事务,自动控制事务提交和回滚
	// using gorm transaction, auto commit or rollback
	return f.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建事务工厂
		// create transaction factory
		txFactory := f.withTransaction(tx)
		// 执行用户逻辑
		// execute user logic
		return fn(txFactory)
	})
}

func (f *repoFactory) User() genericRepo.GenericRepo[model.User, *model.User] {
	return genericRepo.NewGenericRepo[model.User](f.db)
}

func (f *repoFactory) Order() genericRepo.GenericRepo[model.Order, *model.Order] {
	return genericRepo.NewGenericRepo[model.Order](f.db)
}
