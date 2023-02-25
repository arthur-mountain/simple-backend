package databases

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TMysql struct {
	DB             *gorm.DB
	ReconnectCount uint
	DNS            string
	// Others config or gorm.config
}

func MysqlInit(config map[string]interface{}) *TMysql {
	instance := new(TMysql)

	dns, ok := config["DNS"]
	if ok {
		instance.DNS = dns.(string)
	}

	return instance
}

func (t *TMysql) Connect() error {
	connect, err := gorm.Open(mysql.Open(t.DNS), &gorm.Config{})

	if err != nil { // TODO: logger connect error
		fmt.Println("Init database error: ", err)
		return err
	}

	fmt.Println("Database connected")
	t.DB = connect
	return nil
}

func (t *TMysql) WithTrx(trx *gorm.DB) *TMysql {
	t.DB = trx
	return t
}

// TODO: If Exec implement, this method whether or not be implemented?
func (t *TMysql) GetInstance() (*gorm.DB, error) {
	if t.DB == nil {
		if err := t.Connect(); err != nil {
			fmt.Println("reconnect mysql database error")
			return nil, errors.New("reconnect mysql database error")
		}

		// reconnectCount := t.ReconnectCount
		// for reconnectCount > 0 {
		// 	if err := t.Connect(); err == nil {
		// 		break
		// 	}

		// 	time.Sleep(1 * time.Second)
		// 	reconnectCount--
		// }

		// if reconnectCount <= 0 {
		// 	return nil, errors.New("reconnect mysql database error")
		// }
	}

	return t.DB, nil
}

// return bool that means is callback be invoked success
func (t *TMysql) Execute(callback func(DB *gorm.DB) error, model interface{}) error {
	if t.DB == nil {
		count := t.ReconnectCount
		for count > 0 {
			if err := t.Connect(); err == nil {
				break
			}

			time.Sleep(1 * time.Second)
			count--
		}

		if count <= 0 {
			// TODO: may should add logger(reconnect mysql database error)
			return errors.New("reconnect mysql database error")
		}
	}

	if model != nil && t.DB.Migrator().HasTable(&model) {
		return callback(t.DB.Model(&model))
	}

	return callback(t.DB)
}
