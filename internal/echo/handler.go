package echo

import (
	"strings"

	"github.com/Emibotz/megu/pkg/loggers"
	milky "github.com/Szzrain/Milky-go-sdk"
)

func Handler(session *milky.Session, event *milky.ReceiveMessage) {
	logger := loggers.NewDefault("Echo")

	// 必须是群聊或私聊消息
	if event.MessageScene != "group" && event.MessageScene != "friend" {
		return
	}

	// 消息必须是纯文本
	if len(event.Segments) != 1 || event.Segments[0].Type() != "text" {
		return
	}

	textMsg, ok := event.Segments[0].(*milky.TextElement)
	if !ok {
		return
	}

	text := textMsg.Text

	// 检测消息前缀
	if strings.HasPrefix(text, ".echo") {

		// 去除消息前缀
		str := strings.TrimLeft(strings.Replace(text, ".echo", "", 1), " \n")

		// 如果需要重复的消息是空字符串，直接返回
		if str == "" {
			return
		}

		// 构建并发送消息
		logger.Infof("Echo: %s", str)

		reply := []milky.IMessageElement{
			&milky.TextElement{
				Text: str,
			},
		}

		switch event.MessageScene {
		case "group":
			_, err := session.SendGroupMessage(event.PeerId, &reply)
			if err != nil {
				logger.Errorf("Failed to send group message: %v", err)
				return
			}
		case "friend":
			_, err := session.SendPrivateMessage(event.PeerId, &reply)
			if err != nil {
				logger.Errorf("Failed tot send group message: %v", err)
				return
			}
		}
	}
}
