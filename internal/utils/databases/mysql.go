package databases

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	errorConstants "simple-backend/internal/constants/error"
	errorUtils "simple-backend/internal/utils/error"

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

	if err != nil {
		fmt.Println("Init database error: ", err)
		return err
	}

	fmt.Println("Database connected")
	if os.Getenv("ENVIRONMENT") == "production" {
		t.DB = connect
	} else {
		t.DB = connect.Debug()
	}

	return nil
}

func (t *TMysql) WithTrx(trx *gorm.DB) *TMysql {
	t.DB = trx
	return t
}

func (t *TMysql) checkDbIsExistsAndReConnect() error {
	if t.DB != nil {
		return nil
	}

	count := t.ReconnectCount
	for t.Connect() != nil {
		if count <= 0 {
			return errors.New("reconnect mysql database error")
		}

		time.Sleep(1 * time.Second)
		count--
	}

	return nil
}

// could get instance of db, but better wey is use Execute method
func (t *TMysql) GetInstance() (*gorm.DB, error) {
	if err := t.checkDbIsExistsAndReConnect(); err != nil {
		return nil, err
	}

	return t.DB, nil
}

// return error check callback was invoked success
func (t *TMysql) Execute(callback func(DB *gorm.DB) error, model interface{}) *errorUtils.CustomError {
	instance, err := t.GetInstance()

	if err != nil {
		return errorUtils.NewCustomError(err, http.StatusInternalServerError).SetCode(errorConstants.MySQLConnectError)
	}

	if model != nil {
		instance = instance.Model(&model)
	}

	if err = callback(instance); err != nil {
		return errorUtils.CheckRepoError(err)
	}

	return nil
}
