package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Emibotz/megu/internal/echo"
	"github.com/Emibotz/megu/internal/store/pgsql"
	"github.com/Emibotz/megu/internal/todo"
	"github.com/Emibotz/megu/pkg/loggers"
	milky "github.com/Szzrain/Milky-go-sdk"
	"github.com/joho/godotenv"
	"golang.org/x/term"
)

const (
	Title = `
    __  _________________  ______  ____  ______
   /  |/  / ____/ ____/ / / / __ )/ __ \/_  __/
  / /|_/ / __/ / / __/ / / / __  / / / / / /
 / /  / / /___/ /_/ / /_/ / /_/ / /_/ / / /
/_/  /_/_____/\____/\____/_____/\____/ /_/

	`
)

func main() {
	// 创建主函数日志器
	mainLogger := loggers.NewDefault("Main")

	// 创建根上下文
	ctx := context.Background()

	// 打印标题
	width, height, err := term.GetSize(0)
	if err != nil {
		mainLogger.Errorf("Failed to get terminal size: %v", err)
	}

	if width >= 47 && height >= 6 {
		mainLogger.Infof("Welcome to\n%s", Title)
	} else {
		mainLogger.Info("Welcome to MEGUBOT")
	}

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

	// 创建仓库
	pgDB, err := pgsql.New(ctx, os.Getenv("CONN_STRING"))
	if err != nil {
		panic(err)
	}

	todoStore, err := pgDB.TODOStore()
	if err != nil {
		panic(err)
	}

	// 创建待办服务
	todoService := todo.NewService(todoStore)

	// 添加测试事件处理器
	session.AddHandler(func(session *milky.Session, event *milky.ReceiveMessage) {
		logger := loggers.NewDefault("TestHandler1")

		logger.Infof("MessageScene: %s", event.MessageScene)
	})

	// 添加回声处理器
	session.AddHandler(echo.Handler)

	// 添加待办处理器
	todoHandler := todo.NewHandler(ctx, todoService)
	session.AddHandler(todoHandler.Handle)

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
