package repo

import (
	"context"
	"go-tricks/gorm/decoupling-method/model"
	genericRepo "go-tricks/gorm/decoupling-method/repo/generic"
)

type RepoFactory interface {
	Transaction(ctx context.Context, fn func(factory RepoFactory) error) error
	User() genericRepo.GenericRepo[model.User, *model.User]
	Order() genericRepo.GenericRepo[model.Order, *model.Order]
}
