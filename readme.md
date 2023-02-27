### Golang Restful Api


##### Folder Structure
1. config: 設定檔
2. internal/interactor: 共用 struct or middleware
3. internal/domain: 定義主要規範(合同) struct、interface
3. internal/utils: 共用函式
4. migration: sql script

##### Development
開發環境主要使用 [docker-compose](https://www.docker.com/)，

程式語言: [Golang](https://go.dev/)，

其他工具:
1. [Gin](https://github.com/gin-gonic/gin): Web service
2. [Gorm](https://github.com/go-gorm/gorm): Relational mapping database
3. [Mysql](https://www.mysql.com/): Main relational database
4. [Redis](https://redis.io/): Cache database
5. [Golang-jwt](https://github.com/golang-jwt/jwt): JWT token for authentication
6. [Logrus](https://github.com/sirupsen/logrus): Logger for gin http api log
7. [Uuid](https://github.com/google/uuid): Random unique identify id
8. [Swagger](https://github.com/swaggo/swag): API document
9. [Crypto](https://pkg.go.dev/golang.org/x/crypto): Crypto the secret data to database

---

Not implement:
1. [Rebbitmq](https://www.rabbitmq.com/): Message queue service
