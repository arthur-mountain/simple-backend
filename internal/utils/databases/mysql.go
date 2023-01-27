package databases

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlConfig struct {
	DNS *string
	// Others config or gorm.config
}

func (c *MysqlConfig) Connect() (*gorm.DB, error) {
	connect, err := gorm.Open(mysql.Open(*c.DNS), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	fmt.Println("Database connected")
	return connect.Debug(), nil
}
