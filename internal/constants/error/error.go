package constants

type ErrorCode int64

// 自定義錯誤碼mapping到錯誤訊息
const (
	// 服務端錯誤
	ServerError               ErrorCode = 10101
	TooManyRequests           ErrorCode = 10102
	AuthorizationError        ErrorCode = 10103
	AuthorizationExpiredError ErrorCode = 10104
	ParamBindError            ErrorCode = 10105
	CallHTTPError             ErrorCode = 10106
	ConfigGoVersionError      ErrorCode = 20107

	// Mysql錯誤
	MySQLInitError    ErrorCode = 20400
	MySQLConnectError ErrorCode = 20401
	MySQLSearchError  ErrorCode = 20402
	// redis錯誤
	RedisConnectError  ErrorCode = 20500
	RedisSearchError   ErrorCode = 20501
	RedisClearError    ErrorCode = 20502
	RedisSearchIsEmpty ErrorCode = 20503

	// 商業邏輯錯誤
	// User
	IllegalUserName ErrorCode = 20101
	UserCreateError ErrorCode = 20102
	UserUpdateError ErrorCode = 20103
	UserSearchError ErrorCode = 20104
	UserDeleteError ErrorCode = 20105
)

// 自定義錯誤碼的錯誤訊息
var ErrorCodeMapping = map[ErrorCode]string{
	// 服務端錯誤
	ServerError:               "Internal Server Error",
	TooManyRequests:           "Too Many Requests",
	ParamBindError:            "Body 參數錯誤",
	AuthorizationError:        "token 認證失敗",
	AuthorizationExpiredError: "token 已過期",
	CallHTTPError:             "調用第三方 API 錯誤",
	ConfigGoVersionError:      "GoVersion錯誤",

	// Mysql錯誤
	MySQLConnectError: "MySQL連接錯誤",
	MySQLInitError:    "MySQL初始化錯誤",
	MySQLSearchError:  "MySQL查詢錯誤",

	// redis錯誤
	RedisConnectError:  "Redis連接錯誤",
	RedisSearchError:   "RedisKey查詢失敗",
	RedisClearError:    "RedisKey清空失敗",
	RedisSearchIsEmpty: "RedisKey不存在",

	// 商業邏輯錯誤
	IllegalUserName: "非法用戶名",
	UserCreateError: "創建用戶失敗",
	UserUpdateError: "更新用戶失敗",
	UserSearchError: "查詢用戶失敗",
	UserDeleteError: "刪除用戶失敗",
}
