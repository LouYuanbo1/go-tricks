package main

import (
	"context"
	"fmt"
	"go-tricks/gorm/decoupling-method/model"
	repo "go-tricks/gorm/decoupling-method/repo/factory"
	orderService "go-tricks/gorm/decoupling-method/service/order"
	userService "go-tricks/gorm/decoupling-method/service/user"
	"go-tricks/gorm/decoupling-method/table"
	"log"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库
	dsn := "host=localhost user=postgres password=your_password dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("InitDB(dbUtils): failed to connect db: %v", err)
	}
	// 创建用户表
	err = table.NewUserTable(db)
	if err != nil {
		log.Fatalf("NewUserTable(db): failed to create user table: %v", err)
	}

	err = table.NewOrderTable(db)
	if err != nil {
		log.Fatalf("NewOrderTable(db): failed to create order table: %v", err)
	}

	repoFactory := repo.NewRepoFactory(db)

	userService := userService.NewUserService(repoFactory)
	orderService := orderService.NewOrderService(repoFactory)

	users := make([]*model.User, 0, 10)
	// 创建用户
	for i := 0; i < 10; i++ {
		user := &model.User{
			Name:  "testuser" + strconv.Itoa(i),
			Email: "testuser" + strconv.Itoa(i) + "@example.com",
		}
		users = append(users, user)
	}
	err = userService.CreateUsers(context.Background(), users, 100)
	if err != nil {
		log.Fatalf("create user failed: %v", err)
	}
	// 获取用户
	users2, err := userService.GetUsers(context.Background(), []uint64{1, 2, 3})
	if err != nil {
		log.Fatalf("get user failed: %v", err)
	}
	log.Printf("user: %v", users2)
	for _, user := range users2 {
		log.Printf("user: %v", user)
	}

	// 更新用户
	users2[0].Name = "updateduser0"

	err = userService.UpdateUser(context.Background(), users2[0])
	if err != nil {
		log.Fatalf("update user failed: %v", err)
	}

	users2, err = userService.GetUsers(context.Background(), []uint64{1, 2, 3})
	if err != nil {
		log.Fatalf("get user failed: %v", err)
	}
	log.Printf("user: %v", users2)
	for _, user := range users2 {
		log.Printf("user: %v", user)
	}

	user3, err := orderService.CreateOrderWithUser(context.Background(), 2, 100)
	if err != nil {
		log.Fatalf("create order failed: %v", err)
	}
	fmt.Printf("user: %v", user3)
}
