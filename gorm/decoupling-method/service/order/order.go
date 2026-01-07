package service

import (
	"context"
	"fmt"
	"go-tricks/gorm/decoupling-method/model"
	repo "go-tricks/gorm/decoupling-method/repo/factory"
	"log"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	CreateOrders(ctx context.Context, orders []*model.Order, batchSize int) error
	GetOrder(ctx context.Context, id uint64) (*model.Order, error)
	GetOrders(ctx context.Context, ids []uint64) ([]*model.Order, error)
	GetOrdersByPage(ctx context.Context, page, pageSize int) ([]*model.Order, error)
	GetOrdersByCursor(ctx context.Context, cursor uint64, pageSize int) ([]*model.Order, uint64, bool, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	DeleteOrder(ctx context.Context, id uint64) error
	DeleteOrders(ctx context.Context, ids []uint64) error
	//事务代码
	//transaction
	CreateOrderWithUser(ctx context.Context, userID uint64, orderID uint64) (*model.User, error)
}

type orderService struct {
	repoFactory repo.RepoFactory
}

func NewOrderService(repoFactory repo.RepoFactory) OrderService {
	return &orderService{repoFactory: repoFactory}
}

func (o *orderService) CreateOrder(ctx context.Context, order *model.Order) error {
	return o.repoFactory.Order().Create(ctx, order)
}

func (o *orderService) CreateOrders(ctx context.Context, orders []*model.Order, batchSize int) error {
	return o.repoFactory.Order().CreateInBatches(ctx, orders, batchSize)
}

func (o *orderService) GetOrder(ctx context.Context, id uint64) (*model.Order, error) {
	return o.repoFactory.Order().GetByID(ctx, id)
}

func (o *orderService) GetOrders(ctx context.Context, ids []uint64) ([]*model.Order, error) {
	return o.repoFactory.Order().GetByIDs(ctx, ids)
}

func (o *orderService) GetOrdersByPage(ctx context.Context, page, pageSize int) ([]*model.Order, error) {
	return o.repoFactory.Order().GetByPage(ctx, page, pageSize)
}

func (o *orderService) GetOrdersByCursor(ctx context.Context, cursor uint64, pageSize int) ([]*model.Order, uint64, bool, error) {
	return o.repoFactory.Order().GetByCursor(ctx, cursor, pageSize)
}

func (o *orderService) UpdateOrder(ctx context.Context, order *model.Order) error {
	return o.repoFactory.Order().Update(ctx, order)
}

func (o *orderService) DeleteOrder(ctx context.Context, id uint64) error {
	return o.repoFactory.Order().DeleteByID(ctx, id)
}

func (o *orderService) DeleteOrders(ctx context.Context, ids []uint64) error {
	return o.repoFactory.Order().DeleteByIDs(ctx, ids)
}

func (o *orderService) CreateOrderWithUser(ctx context.Context, userID uint64, orderID uint64) (*model.User, error) {
	//可以通过闭包获取任意结果
	//Any result can be obtained through closures
	var user *model.User
	err := o.repoFactory.Transaction(ctx, func(factory repo.RepoFactory) error {
		// 1. 获取事务版本的仓储
		orderRepo := factory.Order()
		userRepo := factory.User()

		// 2. 执行数据库操作（在事务中）
		order := &model.Order{
			UserID:    userID,
			ProductID: orderID,
		}

		var err error

		user, err = userRepo.GetByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		log.Printf("user: %v", user)

		if err := orderRepo.Create(ctx, order); err != nil {
			return fmt.Errorf("create order failed: %w", err)
		}

		log.Printf("order: %v", order)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}
	return user, nil
}
