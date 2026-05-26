package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Emibotz/megu/internal/echo"
	"github.com/Emibotz/megu/pkg/loggers"
	milky "github.com/Szzrain/Milky-go-sdk"
	"github.com/joho/godotenv"
)

func main() {

	// 调试信息
	fmt.Println("Hello, Milky!")

	// 创建主函数日志器
	mainLogger := loggers.NewDefault("Main")

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	wsGateway := os.Getenv("WS_ADDR")
	restGateway := os.Getenv("HTTP_ADDR")
	token := os.Getenv("TOKEN")

	// 创建会话
	sessionLogger := loggers.NewDefault("Session")

	session, err := milky.New(wsGateway, restGateway, token, sessionLogger)
	if err != nil {
		panic(err)
	}

	// 添加测试事件处理器
	session.AddHandler(func(session *milky.Session, event *milky.ReceiveMessage) {
		logger := loggers.NewDefault("TestHandler1")

		logger.Infof("MessageScene: %s", event.MessageScene)
	})

	// 添加回声处理器
	session.AddHandler(echo.Handler)

	// 开启会话
	if err := session.Open(); err != nil {
		panic(err)
	}

	// Graceful Terminate
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	mainLogger.Info("Graceful Terminating...")

	if err := session.Close(); err != nil {
		mainLogger.Errorf("Err terminating program: %v", err)
	}

	mainLogger.Info("The program should be terminated.")
}
