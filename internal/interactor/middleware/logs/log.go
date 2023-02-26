package logs

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	errorUtils "simple-backend/internal/utils/error"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type tLogHook struct{}

// Notice: docker container always +0UTC
func (h *tLogHook) Fire(entry *logrus.Entry) error {
	timezoneStr := os.Getenv("TIMEZONE")
	timezone, _ := time.LoadLocation(timezoneStr)

	// set log timezone
	entry.Time = entry.Time.In(timezone)
	entry.Data["timezone"] = timezoneStr
	// set log formatter
	entry.Logger.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@time",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
		},
	}

	return nil
}

func (h *tLogHook) Levels() []logrus.Level {
	if os.Getenv("ENVIRONMENT") == "production" {
		return []logrus.Level{
			logrus.InfoLevel,
			logrus.WarnLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}
	}

	return logrus.AllLevels
}

func NewLogger() *logrus.Logger {
	// get log path
	logPath, err := os.Getwd()
	if err != nil {
		fmt.Println("new logger couldn't found current directory", err)
	}

	// create log directory
	logPath = logPath + "/logs/"
	err = os.MkdirAll(logPath, 0777)
	if err != nil {
		fmt.Println("create log directory", err)
	}

	// check log file path is exists
	fileName := path.Join(
		logPath,                                // log path
		time.Now().Format("2006-01-02")+".log", // log filename
	)
	if _, err = os.Stat(fileName); err != nil {
		if _, err = os.Create(fileName); err != nil {
			fmt.Println("create file error: ", err)
		}
	}

	// open log file
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("open log file err, file is not exists", err)
	}

	// create logrus instance
	logger := logrus.New()
	logger.SetOutput(file)        // append log to file
	logrus.SetReportCaller(true)  // 測試環境，將檔案名稱和位置一同紀錄
	logger.AddHook(new(tLogHook)) // invoked before append log to file

	return logger
}

func handleRecoverLog(c *gin.Context, logger *logrus.Logger) {
	if err := recover(); err != nil {
		customError := errorUtils.NewCustomError(
			errors.New("please contact the administrator"),
			http.StatusInternalServerError,
		)
		c.AbortWithStatusJSON(customError.HttpStatusCode, customError)

		statusCode := c.Writer.Status()
		reqUri := c.Request.RequestURI
		reqMethod := c.Request.Method
		clientIP := c.ClientIP()
		logger.Errorf("code:%d | ip:%s | method:%s | uri:%s | message: %s",
			statusCode,
			clientIP,
			reqMethod,
			reqUri,
			err,
		)
	}
}

func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := NewLogger()
		defer handleRecoverLog(c, logger)
		c.Next()

		if httpCode := c.Writer.Status(); httpCode >= 400 {
			reqUri := c.Request.RequestURI
			reqMethod := c.Request.Method
			clientIP := c.ClientIP()
			logger.Errorf("code:%d | ip:%s | method:%s | uri:%s",
				httpCode,
				clientIP,
				reqMethod,
				reqUri,
			)
		}
	}
}
