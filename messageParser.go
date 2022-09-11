package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Parses a full payload from the server, which may include multiple messages
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
			if roomId != "" {
				a.parseRoomServerMessage(roomId, payloadLines[i])
			} else {
				a.parseGlobalServerMessage(payloadLines[i])
			}
		} else {
			goPrint("blank line in message payload")
		}

	}
}

// Parse one message at a time
func (a *App) parseGlobalServerMessage(message string) {
	// parse a message from the server
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
		goPrint("todo: handle query response", message)
	default:
		goPrint("<!>WARNING<!> unhandled global message", message)
	}
}

// messages sent in the context of a room
func (a *App) parseRoomServerMessage(roomId, message string) {
	// parse a message from the server
	chunkedMsg := NewSplitServerMessage(message)
	msgType := MessageType(strings.ToLower(chunkedMsg.Get(1)))
	// need to handle ||MESSAGE (two pipes) and MESSAGE (no pipes)
	if message[0] != '|' {
		goPrint("Just a message", message)
	}
	switch msgType {
	case Init:
		//   |init|<room type>
		// 0 | 1  |     2
		payload := NewRoomPayload{roomId, RoomType(chunkedMsg.Get(2))}
		a.channels.frontendChan <- ShowdownEvent{NewRoomTopic, payload}
	case DeInit:
		//   |deinit
		a.channels.frontendChan <- ShowdownEvent{RoomExitTopic, roomId}
	case Title:
		//   |title|<title>
		// 0 |  1  |   2
		//payload := RoomStatePayload{Title: chunkedMsg.Get(2)}
		payload := RoomStatePayload{RoomId: roomId, Title: chunkedMsg.Get(2)}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Users:
		//   |users|<user list>
		// 0 |  1  |     2
		// user list is comma delimited
		// @ delimits a user entry with their status message
		// status messages starting with ! means the user is aways
		userList := NewSplitString(chunkedMsg.ReassembleTail(2), ",")
		users := make([]string, len(userList.inner))
		for i, u := range userList.inner {
			users[i] = NewUser(u).UserName
		}
		payload := RoomStatePayload{RoomId: roomId, Users: users}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Html:
		//   |html|<html>
		// 0 | 1  |  2
		payload := RoomHtmlPayload{roomId, chunkedMsg.Get(2), "", false}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case Uhtml:
		//   |uhtml|<name>|<html>
		// 0 |  1  |  2   |  3
		// uhtml contains named, updatable html
		payload := RoomHtmlPayload{roomId, chunkedMsg.Get(3), chunkedMsg.Get(2), false}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case UhtmlChange:
		//   |uhtmlchange|<name>|<html>
		// 0 |     1     |  2   |  3
		// uhtmlchange updates the named updatable html
		payload := RoomHtmlPayload{roomId, chunkedMsg.Get(3), chunkedMsg.Get(2), true}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case Join, Join2:
		//   |join|<user>
		// 0 | 1  |  2
		joinUser := NewUser(chunkedMsg.Get(2))
		payload := RoomMessagePayload{roomId, "system", fmt.Sprintf("%c%s joined", joinUser.Rank, joinUser.UserName)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case Leave, Leave2:
		//   |leave|<user>
		// 0 |  1  |  2
		leaveUser := NewUser(chunkedMsg.Get(2))
		payload := RoomMessagePayload{roomId, "system", fmt.Sprintf("%s joined", leaveUser.UserName)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case Name, Name2:
		//   |name|<new user name>|<old id>
		// 0 | 1  |       2       |   3
		nameUser := NewUser(chunkedMsg.Get(2))
		payload := RoomMessagePayload{roomId, "system", fmt.Sprintf("%s changed their name to %s", chunkedMsg.Get(3), nameUser.UserName)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case ChatMsg, ChatMsg2:
		//   |chat|<user>|<message>
		// 0 | 1  |  2   |    3
		chatUser := NewUser(chunkedMsg.Get(2))
		payload := RoomMessagePayload{roomId, chatUser.UserName, chunkedMsg.ReassembleTail(3)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case Notify:
		//   |notify|<title>|<message - optional>|<highlight token - optional
		// 0 |  1   |   2   |         3          |             4
		goPrint("todo: handle notify", message)
	case Timestamp2:
		//   |:|<unix timestamp>
		goPrint("ignoring server room timestamp: ", message)
	case ChatTs:
		//   |c:|<unix timestamp>|<user>|<message>
		// 0 |1 |        2       |  3   |    4
		// "The exact fate of this command is uncertain - it may or may not be replaced with a more generalized way to transmit timestamps at some point." - sounds like this might not even be in use
		goPrint("todo: handle timestamped room chat", message)
	case BattleStart, BattleStart2:
		//   |battle|<room id>|<user 1>|<user 2>
		// 0 |  1   |    2    |   3    |   4
		u1 := NewUser(chunkedMsg.Get(3))
		u2 := NewUser(chunkedMsg.Get(4))
		payload := RoomMessagePayload{roomId, "system", fmt.Sprintf("a battle has started between %s and %s in %s", u1.UserName, u2.UserName, chunkedMsg.Get(2))}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	default:
		a.parseBattleMessage(roomId, chunkedMsg)
	}
}

func NewSplitServerMessage(msg string) *SplitString {
	return NewSplitString(msg, "|")
}
