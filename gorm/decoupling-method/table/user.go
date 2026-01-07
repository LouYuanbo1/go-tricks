package table

import (
	"log"

	"gorm.io/gorm"
)

func NewUserTable(db *gorm.DB) error {
	table := `CREATE TABLE IF NOT EXISTS users(
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(50) NOT NULL,
			email VARCHAR(50) NOT NULL
		);`
	err := db.Exec(table).Error
	if err != nil {
		log.Printf("NewUserTable(table): 创建用户表失败: %v", err)
		return err
	}
	return nil
}
