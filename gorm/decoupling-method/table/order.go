package table

import (
	"log"

	"gorm.io/gorm"
)

func NewOrderTable(db *gorm.DB) error {
	table := `CREATE TABLE IF NOT EXISTS orders(
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL
		);`
	err := db.Exec(table).Error
	if err != nil {
		log.Printf("NewOrderTable(table): 创建订单表失败: %v", err)
		return err
	}
	return nil
}
