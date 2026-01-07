package service

import (
	"context"
	"go-tricks/gorm/decoupling-method/model"
	repo "go-tricks/gorm/decoupling-method/repo/factory"
)

type UserService interface {
	CreateUser(ctx context.Context, user *model.User) error
	CreateUsers(ctx context.Context, users []*model.User, batchSize int) error
	GetUser(ctx context.Context, id uint64) (*model.User, error)
	GetUsers(ctx context.Context, ids []uint64) ([]*model.User, error)
	GetUsersByPage(ctx context.Context, page, pageSize int) ([]*model.User, error)
	GetUsersByCursor(ctx context.Context, cursor uint64, pageSize int) ([]*model.User, uint64, bool, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id uint64) error
	DeleteUsers(ctx context.Context, ids []uint64) error
}

type userService struct {
	repoFactory repo.RepoFactory
}

func NewUserService(repoFactory repo.RepoFactory) UserService {
	return &userService{repoFactory: repoFactory}
}

func (u *userService) CreateUser(ctx context.Context, user *model.User) error {
	return u.repoFactory.User().Create(ctx, user)
}

func (u *userService) CreateUsers(ctx context.Context, users []*model.User, batchSize int) error {
	return u.repoFactory.User().CreateInBatches(ctx, users, batchSize)
}

func (u *userService) GetUser(ctx context.Context, id uint64) (*model.User, error) {
	return u.repoFactory.User().GetByID(ctx, id)
}

func (u *userService) GetUsers(ctx context.Context, ids []uint64) ([]*model.User, error) {
	return u.repoFactory.User().GetByIDs(ctx, ids)
}

func (u *userService) GetUsersByPage(ctx context.Context, page, pageSize int) ([]*model.User, error) {
	return u.repoFactory.User().GetByPage(ctx, page, pageSize)
}

func (u *userService) GetUsersByCursor(ctx context.Context, cursor uint64, pageSize int) ([]*model.User, uint64, bool, error) {
	return u.repoFactory.User().GetByCursor(ctx, cursor, pageSize)
}

func (u *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return u.repoFactory.User().Update(ctx, user)
}

func (u *userService) DeleteUser(ctx context.Context, id uint64) error {
	return u.repoFactory.User().DeleteByID(ctx, id)
}

func (u *userService) DeleteUsers(ctx context.Context, ids []uint64) error {
	return u.repoFactory.User().DeleteByIDs(ctx, ids)
}
