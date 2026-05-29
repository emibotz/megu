package todo

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Emibotz/megu/pkg/loggers"
	"github.com/Emibotz/megu/pkg/utils"
	milky "github.com/Szzrain/Milky-go-sdk"
	"github.com/go-andiamo/splitter"
)

type handler struct {
	logger milky.Logger

	ctx context.Context

	todoService *Service
}

func NewHandler(
	ctx context.Context,

	todoService *Service,
) *handler {
	return &handler{
		logger:      loggers.NewDefault("TODOHandler"),
		ctx:         ctx,
		todoService: todoService,
	}
}

func (h *handler) list(ctx context.Context, page int) (string, error) {
	list, err := h.todoService.List(ctx, page)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}

	for _, listLine := range list {
		completed := " "
		if listLine.TD.Completed {
			completed = "+"
		}

		fmt.Fprintf(&builder, "\n[%d] [%s] %s", listLine.Num, completed, listLine.TD.Title)
	}

	return builder.String()[1:], nil
}

func (h *handler) Handle(session *milky.Session, event *milky.ReceiveMessage) {
	// 消息必须是纯文本
	if len(event.Segments) != 1 || event.Segments[0].Type() != milky.Text {
		return
	}

	textSeg, ok := event.Segments[0].(*milky.TextElement)
	if !ok {
		return
	}

	text := textSeg.Text

	// 把消息用空格分隔开，开始解析语法
	split, err := splitter.NewSplitter(' ', splitter.DoubleQuotesBackSlashEscaped)
	if err != nil {
		utils.SendTextMessage(session, event.MessageScene, event.PeerId, "发生致命错误，请告知管理员。todo/handler.go:72")
		return
	}

	args, err := split.Split(text, splitter.IgnoreEmpties)
	if err != nil {
		return
	}

	// 移除首尾处的引号，把文本内部的转义引号 `\"` 转换成普通引号 `"`
	for i, arg := range args {
		args[i] = strings.ReplaceAll(strings.TrimPrefix(strings.TrimSuffix(arg, "\""), "\""), "\\\"", "\"")
	}

	if args[0] != ".todo" && args[0] != ".td" {
		return
	}

	length := len(args)

	// 如果没有参数，返回列表
	if length == 1 {
		reply, err := h.list(h.ctx, 0)
		if err != nil {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Failed to list TODOs: %v", err))
			return
		}

		utils.SendTextMessage(session, event.MessageScene, event.PeerId, reply)
		return
	}

	switch args[1] {
	// 创建
	case "new":

		// 检测参数数量
		if length < 3 || length > 4 {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Expected 3 or 4 arguments, got %d", length))
			return
		}

		// 读取参数
		title := args[2]
		content := ""

		if length == 4 {
			content = args[3]
		}

		// 创建待办
		if err := h.todoService.Create(context.TODO(), title, content); err != nil {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Failed to create TODO: %v", err))
			return
		}

		// 调试信息
		h.logger.Infof("User %d Created TODO", event.SenderId)

		utils.SendTextMessage(session, event.MessageScene, event.PeerId, "Created!")

		// 列表
	case "list":

		// 检测参数数量
		if length > 3 {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Too many arguments! Expected 2 or 3 arguments, got %d", length))
			return
		}

		// 计算页数
		page := 0
		if length == 3 {
			p, err := strconv.Atoi(args[2])
			if err != nil || p < 1 {
				utils.SendTextMessage(session, event.MessageScene, event.PeerId, "请输入一个正确的正整数！")
			}

			page = p - 1
		}

		// 获取分页内容
		list, err := h.list(h.ctx, page)
		if err != nil {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Failed to list TODOs: %v", err))
			return
		}

		// 调试信息
		h.logger.Infof("User %d Listed Page %d", event.SenderId, page)

		utils.SendTextMessage(session, event.MessageScene, event.PeerId, list)

	default:

		// 检测参数数量
		if length > 2 {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Invalid argument: %s", args[1]))
		}

		n, err := strconv.Atoi(args[1])
		if err != nil {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Invalid argument: %s", args[1]))
			return
		}

		// 获取待办信息
		td, err := h.todoService.GetByNumber(h.ctx, n)
		if err != nil {
			utils.SendTextMessage(session, event.MessageScene, event.PeerId, fmt.Sprintf("Failed to get: %v", err))
			return
		}

		// 构建回复消息
		builder := strings.Builder{}

		fmt.Fprintf(&builder, "Title: %s", td.Title)

		if td.Content != "" {
			fmt.Fprintf(&builder, "\nContent: %s", td.Content)
		}

		fmt.Fprintf(&builder, "\nCreated At: %s", td.CreatedAt.Format(time.DateTime))

		if td.DeletedAt != nil {
			fmt.Fprintf(&builder, "\nDeleted At: %s", td.DeletedAt.Format(time.DateTime))
		}

		if td.Completed {
			fmt.Fprint(&builder, "\nCompleted: TRUE")
		} else {
			fmt.Fprint(&builder, "\nCompleted: FALSE")
		}

		utils.SendTextMessage(session, event.MessageScene, event.PeerId, builder.String())

	}
}
