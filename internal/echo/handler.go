package echo

import (
	"strings"

	"github.com/Emibotz/megu/pkg/loggers"
	"github.com/Emibotz/megu/pkg/utils"
	milky "github.com/Szzrain/Milky-go-sdk"
)

// 回声命令处理器
func Handler(session *milky.Session, event *milky.ReceiveMessage) {
	logger := loggers.NewDefault("Echo")

	// 消息必须是纯文本
	if len(event.Segments) != 1 || event.Segments[0].Type() != milky.Text {
		return
	}

	textSeg, ok := event.Segments[0].(*milky.TextElement)
	if !ok {
		return
	}

	text := textSeg.Text

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

		_, err := utils.SendTextMessage(session, event.MessageScene, event.PeerId, str)
		if err != nil {
			logger.Errorf("Failed to send message: %v", err)
		}
	}
}
