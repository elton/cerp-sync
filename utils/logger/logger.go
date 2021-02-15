package logger

import (
	"io"
	"log"
	"os"
)

var (
	// Info info level log for logger.
	Info *log.Logger
	// Error error log for logger
	Error *log.Logger
)

func init() {
	//日志输出文件
	file, err := os.OpenFile("cerp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Faild to open error logger file:", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	//自定义日志格式
	Info = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	log.SetOutput(mw)
}
