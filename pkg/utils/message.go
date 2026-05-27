package utils

import (
	"fmt"

	milky "github.com/Szzrain/Milky-go-sdk"
)

func SendMessage(session *milky.Session, messageScene string, peerID int64, message *[]milky.IMessageElement) (*milky.MessageRet, error) {
	switch messageScene {
	case "friend":
		return session.SendPrivateMessage(peerID, message)
	case "group":
		return session.SendGroupMessage(peerID, message)
	default:
		return nil, fmt.Errorf("unavailable message scene: %s", messageScene)
	}
}

func SendTextMessage(session *milky.Session, messageScene string, peerID int64, text string) (*milky.MessageRet, error) {
	message := []milky.IMessageElement{
		&milky.TextElement{
			Text: text,
		},
	}

	return SendMessage(session, messageScene, peerID, &message)
}
