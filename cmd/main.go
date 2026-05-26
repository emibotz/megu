package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Emibotz/megu/pkg/loggers"
	milky "github.com/Szzrain/Milky-go-sdk"
	"github.com/joho/godotenv"
)

func main() {
	// 调试信息
	fmt.Println("Hello, Milky!")

	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	wsGateway := os.Getenv("WS_ADDR")
	restGateway := os.Getenv("HTTP_ADDR")
	token := os.Getenv("TOKEN")

	// 创建会话
	logger := loggers.NewDefault("MilkyBot")

	session, err := milky.New(wsGateway, restGateway, "Bearer "+token, logger)
	if err != nil {
		panic(err)
	}

	// 添加测试事件处理器
	logger = loggers.NewDefault("TestHandler1")
	session.AddHandler(func(session *milky.Session, event *milky.ReceiveMessage) {
		logger := logger

		logger.Infof("MessageScene: ", event.MessageScene)
	})

	// 开启会话
	if err := session.Open(); err != nil {
		panic(err)
	}

	// Graceful Terminate
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	fmt.Println("Graceful Terminating...")

	if err := session.Close(); err != nil {
		fmt.Printf("Err terminating program: %v", err)
	}

	fmt.Println("The program should be terminated.")
}
