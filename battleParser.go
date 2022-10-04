package main

import (
	"fmt"
	"strconv"
)

// for parsing battle messages from server

// handles battle flow messages and then delegates action messages
func (a *App) parseBattleMessage(roomId string, msg *SplitString) {
	pTrue := true
	pFalse := false
	fullMessage := msg.ReassembleTail(0)
	goPrint("incoming battle message", fullMessage)
	msgType := MessageType(msg.Get(1))
	switch msgType {
	case Request:
		//   |request|<request json>
		// 0 |   1   |      2
		payload := UpdateRoomStatePayload{RoomId: roomId}
		if msg.Get(2) != "" {
			payload.Request = msg.ReassembleTail(2)
			a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
		}
	case Timestamp:
		//   |t:|<unix timestamp>
		// 0 |1 |        2
		goPrint("ignoring battle room timestamp", fullMessage)
	case GameType:
		//   |gametype|<game type>
		// 0 |   1    |     2
		payload := UpdateRoomStatePayload{RoomId: roomId, GameType: msg.ReassembleTail(2)}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Player:
		//   |player|<id>|<name>|<avatar name>|<rating>
		// 0 |  1   | 2  |  3   |      4      |   5
		subPayload := UpdatePlayerPayload{PlayerId: msg.Get(2), Name: msg.Get(3), Avatar: msg.Get(4), Rating: msg.Get(5)}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: &subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case TeamSize:
		//   |teamsize|<id>|<size>
		// 0 |   1    | 2  |  3
		s, _ := strconv.Atoi(msg.Get(3))
		subPayload := UpdatePlayerPayload{PlayerId: msg.Get(2), TeamSize: s}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: &subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Generation:
		//   |gen|<number>
		// 0 | 1 |   2
		g, _ := strconv.Atoi(msg.Get(2))
		payload := UpdateRoomStatePayload{RoomId: roomId, Gen: g}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Tier:
		//   |tier|<format name>
		// 0 | 1  |      2
		payload := UpdateRoomStatePayload{RoomId: roomId, Tier: msg.ReassembleTail(2)}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Rule:
		//   |rule|<rule>
		// 0 | 1  |  2
		// rules are delimited with a colon between name and description
		ruleSplit := NewSplitString(msg.ReassembleTail(2), ":")
		payload := RoomMessagePayload{roomId, ruleSplit.Get(0), ruleSplit.ReassembleTail(1)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case IsRatedBattle:
		//   |rated|<message - optional>
		payload := UpdateRoomStatePayload{RoomId: roomId, Rated: true}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
		if msg.Get(2) != "" {
			payload2 := RoomMessagePayload{roomId, "system", msg.ReassembleTail(2)}
			a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload2}
		}
	case ClearPoke:
		//   |clearpoke
		// "marks the start of team preview"
		goPrint(fullMessage, "is a team preview feature")
	case PreviewPoke:
		//   |poke|<id>|<details>|<item>
		// 0 | 1  | 2  |    3    |  4
		// id is the player id, eg p1
		goPrint(fullMessage, "is a team preview feature")
	case TeamPreview:
		//   |teampreview
		// signals the end of |poke| messages
		goPrint(fullMessage, "is a team preview feature")
	case SimStart:
		//   |start
		sPayload := UpdateRoomStatePayload{RoomId: roomId, Active: &pTrue}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, sPayload}
		cPayload := RoomMessagePayload{roomId, "system", "The battle has begun!"}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, cPayload}
	case Turn:
		//   |turn|<number>
		// 0 | 1  |   2
		payload := RoomMessagePayload{roomId, "turn", msg.ReassembleTail(2)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, payload}
	case TimerOn:
		//   |inactive|<message>
		// 0 |    2   |    3
		sPayload := UpdateRoomStatePayload{RoomId: roomId, Timer: &pTrue}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, sPayload}
		cPayload := RoomMessagePayload{roomId, "timer", msg.ReassembleTail(3)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, cPayload}
	case TimerOff:
		//   |inactiveoff|<message>
		// 0 |     2     |    3
		sPayload := UpdateRoomStatePayload{RoomId: roomId, Timer: &pFalse}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, sPayload}
		cPayload := RoomMessagePayload{roomId, "timer", msg.ReassembleTail(3)}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, cPayload}
	case Upkeep:
		//   |upkeep
	case Win:
		//   |win|<user>
		// 0 | 1 |  2
		sPayload := UpdateRoomStatePayload{RoomId: roomId, Active: &pFalse}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, sPayload}
		cPayload := RoomMessagePayload{roomId, "system", fmt.Sprintf("%s won!", msg.ReassembleTail(2))}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, cPayload}
	case Tie:
		//   |tie
		sPayload := UpdateRoomStatePayload{RoomId: roomId, Active: &pFalse}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, sPayload}
		cPayload := RoomMessagePayload{roomId, "system", "battle ended in a tie"}
		a.channels.frontendChan <- ShowdownEvent{RoomMessageTopic, cPayload}

	default:
		a.parseMajorBattleAction(roomId, msg)
		a.parseMinorBattleAction(roomId, msg)
	}
}

func (a *App) parseMajorBattleAction(roomId string, msg *SplitString) {
	fullMessage := msg.ReassembleTail(0)
	goPrint("possible major battle action", fullMessage)
	msgType := MessageType(msg.Get(1))
	switch msgType {
	case Move:
		//
	case Switch, Drag, Replace, DetailsChanged, FormeChange:
		//   |<type>|<position spec>|<details spec>|<hp spec>
		// 0 |  1   |       2       |       3      |    4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		d := NewPokemonDetails(msg.Get(3))
		h := NewHPStatus(msg.Get(4))
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = &UpdatePlayerPokemon{msgType, p, d, h}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
		// TODO send some sort of room message
	case Swap:
		//
	case Cannot:
		//
	case Faint:
		//   |faint|<position spec>
		// 0 |  1  |       2
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = &UpdatePlayerPokemon{Reason: msgType, Position: p}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	default:
		goPrint("not a supported major Battle action")
	}
}

func (a *App) parseMinorBattleAction(roomId string, msg *SplitString) {
	fullMessage := msg.ReassembleTail(0)
	goPrint("possible minor battle action", fullMessage)
	msgType := MessageType(msg.Get(1))
	switch msgType {
	case Fail:
		//
	case Block:
		//
	case NoTarget:
		//
	case Miss:
		//
	case Damage, Heal:
		//   |<type>|<position spec>|<hp spec>|<from details
		// 0 |  1   |       2       |    3    |      4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		h := NewHPStatus(msg.Get(3))
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = &UpdatePlayerPokemon{Reason: "hpupdate", Position: p, HP: h}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: subPayload}
		goPrint("hp update room state payload", fmt.Sprintf("%+v", payload), fmt.Sprintf("%+v", subPayload), fmt.Sprintf("%+v", h))
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
		// TODO: send a chat message
	default:
		goPrint("not a supported minor battle action", msg.ReassembleTail(0))
	}
}
