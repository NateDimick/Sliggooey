package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

func (a *App) parseServerPayload(payload string) {
	// messages come from the server grouped by an optional roomId
	// unpack the roomId and process each message
	payloadLines := strings.Split(payload, "\n")
	msgStartIndex := 0
	roomId := ""
	if payloadLines[0][0] == '>' {
		goPrint("Message payload for roomId:", payloadLines[0][1:])
		msgStartIndex = 1
		roomId = payloadLines[0][1:]
	} else {
		goPrint("Message payload for the lobby or global room")
	}
	for i := msgStartIndex; i < len(payloadLines); i++ {
		if len(payloadLines[i]) > 0 {
			a.parseServerMessage(roomId, payloadLines[i])
		} else {
			goPrint("blank line in message payload")
		}

	}
}

func (a *App) parseServerMessage(roomId string, message string) {
	// parse a message from the server
	if message[0] != '|' {
		goPrint("Just a message", message)
	}
	chunkedMsg := NewSplitServerMessage(message)
	msgType := MessageType(chunkedMsg.Get(1))
	switch msgType {
	// in order as they are sent from the server, and then in order of
	// global messages, room initialization, room messages, battle messages
	// tournament messages not supported yet
	// for convenience, message format is included in each case
	case ChallStr:
		//   |challstr|<challstr>
		// 0 |    1   |     2
		goPrint("storing challstr", message)
		a.conn.challengeString = chunkedMsg.ReassembleTail(2)
	case Formats:
		//   |formats|format list...
		// 0 |   1   |...
		a.formatList(chunkedMsg.ReassembleTail(2))
	case UpdateSearch:
		//   |updatesearch|<search json>
		// 0 |     1      |     2
		goPrint("update search", message)
		a.updateSearchEvent(chunkedMsg.ReassembleTail(2))
	case UpdateUser:
		//   |updateuser|username|guestflag|avatarnum|<settings json>
		// 0 |     1    |    2   |    3    |    4    |       5
		// guestflag is 0 for guest, 1 for logged in user
		goPrint("update user", message)
		flag, _ := strconv.Atoi(chunkedMsg.Get(3))
		avatar, _ := strconv.Atoi(chunkedMsg.Get(4))
		settings := new(UserSettings)
		err := json.Unmarshal([]byte(chunkedMsg.ReassembleTail(5)), settings)
		if err != nil {
			goPrint("could not parse user settings")
			return
		}
		a.updateUserResultEvent(chunkedMsg.Get(2), flag, avatar, settings)
	case UpdateChallenges:
		//   |updatechallenges|<challenge json>
		// 0 |        1       |       2
		a.updateChallengeEvent(chunkedMsg.ReassembleTail(2))
	case Popup:
		//   |popup|<message>
		// 0 |  1  |    2
		a.popupMessageEvent(chunkedMsg.ReassembleTail(2))
	case UserCount:
		//   |usercount|<count>
		// 0 |    1    |   2
		count, err := strconv.Atoi(chunkedMsg.Get(2))
		if err != nil {
			defer recover()
			goPrint("usercount could not parse number", message)
			panic(err)
		}
		a.userCountEvent(count)
	case NameTaken:
		//   |nametaken|<username>|<reason>
		// 0 |    1    |     2    |    3
		a.nameTakenEvent(chunkedMsg.Get(2), chunkedMsg.ReassembleTail(3))
	case PrivateMessage:
		//   |pm|<sender>|<receiver>|<message>
		// 0 | 1|    2   |     3    |    4
		goPrint("full pm", message)
		if chunkedMsg.Get(4)[0] == '/' {
			goPrint("is a chat command")
			a.parseServerCommand(chunkedMsg.Get(2), chunkedMsg.Get(3), chunkedMsg.ReassembleTail(4))
		} else {
			a.privateMessageEvent(chunkedMsg.Get(2), chunkedMsg.Get(3), chunkedMsg.ReassembleTail(4))
		}
	case QueryResponse:
		//   |queryresponse|<query type>|<json>
		// 0 |      1      |     2      |   3
	default:
		goPrint("TODO: parse message", message)
		if message[0] != '|' {
			// put message directly in to chat with roomId
		} else {
			a.parseBattleMessage(chunkedMsg)
		}
	}
}

func NewSplitServerMessage(msg string) *SplitString {
	return NewSplitString(msg, "|")
}
