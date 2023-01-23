### Golang Restful Api


##### Folder Structure
1. config: 設定檔
2. internal/interactor: 共用 struct or middleware
3. internal/domain: 定義主要規範(合同) struct、interface
3. internal/utils: 共用函式
4. migration: sql script

##### Development
開發環境主要使用 docker-compose，
以下為 docker-compose 有建立起的服務

程式語言: Golang
其他工具:
1. Gin: Golang web
2. Gorm: Golang database
3. Mysql: Main Relation Database
4. Redis: Cache database
5. Rebbitmq: Queue service
6. Golang-jwt: Golang JWT token for authentication
