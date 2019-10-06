package config

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"goweb/util"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func CreateLogFile() *log.Logger {
	//设定日志文件夹 必须在gin之前
	createLogFile("./logs", "/vain.log")
	f, _ := os.Create("./logs/vain.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	config := gin.LoggerConfig{
		Output:    io.MultiWriter(f, os.Stdout),
		SkipPaths: []string{"/test"},
		Formatter: func(params gin.LogFormatterParams) string {
			return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				params.ClientIP,
				params.TimeStamp.Format(time.RFC1123),
				params.Method,
				params.Path,
				params.Request.Proto,
				params.StatusCode,
				params.Latency,
				params.Request.UserAgent(),
				params.ErrorMessage,
			)
		},
	}
	gin.LoggerWithConfig(config)
	logger := log.New(f, "", log.Llongfile)
	return logger
}

//注入log到上线文中
func InjectLog(log *log.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Set("log", log)
		context.Next()
	}
}

var accessChannel = make(chan string, 100)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LogSetUp() gin.HandlerFunc {

	createLogFile("./logs", "/request.log")
	go handleAccessChannel("./logs/request.log")

	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 开始时间
		startTime := util.GetCurrentMilliUnix()

		// 处理请求
		c.Next()

		responseBody := bodyLogWriter.body.String()

		// 结束时间
		endTime := util.GetCurrentMilliUnix()

		if c.Request.Method == "POST" {
			_ = c.Request.ParseForm()
		}

		accessChannel <-
			fmt.Sprintf("[%s] - %s | %s | %s | %s | \"%s %s %s  %s %s \"\n",
				util.GetCurrentDate(),
				c.ClientIP(),
				c.Request.RequestURI,
				c.Request.Method,
				c.Request.Proto,
				c.Request.PostForm.Encode(),
				responseBody,
				c.Request.UserAgent(),
				fmt.Sprintf("%vms", endTime-startTime),
				c.Request.Referer(),
			)
	}
}

func createLogFile(dir string, fileName string) bool {
	if _, err := os.Stat(strings.Join([]string{dir, fileName}, "")); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		_ = os.MkdirAll(dir, 0777)
		f, _ := os.Create(strings.Join([]string{dir, fileName}, ""))
		defer f.Close()
		_ = os.Chmod(strings.Join([]string{dir, fileName}, ""), 0777)
	}
	return true
}

func handleAccessChannel(path string) {
	if f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666); err != nil {
		log.Println(err)
	} else {
		for accessLog := range accessChannel {
			_, _ = f.WriteString(accessLog + "\n")
		}
	}
	return
}
