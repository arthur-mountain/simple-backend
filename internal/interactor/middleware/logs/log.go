package logs

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	errorUtils "simple-backend/internal/utils/error"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	now := time.Now() // current time
	logPath := ""     // folder path
	logFileName := now.Format("2006-01-02") + ".log"

	// set up log path and file name
	dir, err := os.Getwd()
	if err == nil {
		logPath = dir + "/logs/"
	}

	err = os.MkdirAll(logPath, 0777)
	if err != nil {
		fmt.Println(err.Error())
	}

	fileName := path.Join(logPath, logFileName)

	// check log file path is exists
	_, err = os.Stat(fileName)
	if err != nil {
		_, err = os.Create(fileName)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// open log file
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("open log file err: ", err)
	}

	// create logrus instance
	logger := logrus.New()
	logger.Out = src
	logger.SetLevel(logrus.DebugLevel)
	// log time is +0
	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger
}

func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := NewLogger()
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		defer func() {
			if err := recover(); err != nil {
				customError := errorUtils.NewCustomError(
					errors.New("please contact the administrator"),
					http.StatusInternalServerError,
				)
				logger.Errorf("code:%d | ip:%s | method:%s | uri:%s | message: %s",
					statusCode,
					clientIP,
					reqMethod,
					reqUri,
					err,
				)
				c.AbortWithStatusJSON(customError.HttpStatusCode, customError)
			}
		}()

		c.Next()
		logger.Infof("code:%d | ip:%s | method:%s | uri:%s",
			statusCode,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
