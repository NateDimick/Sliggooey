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
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case TeamSize:
		//   |teamsize|<id>|<size>
		// 0 |   1    | 2  |  3
		s, _ := strconv.Atoi(msg.Get(3))
		subPayload := UpdatePlayerPayload{PlayerId: msg.Get(2), TeamSize: s}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: subPayload}
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
		a.parseFieldBattleMessage(roomId, msg)
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
		delta := PokeDelta{HP: h}
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Details: d, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
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
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
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
		delta := PokeDelta{HP: h}
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint("hp update room state payload", fmt.Sprintf("%+v", payload), fmt.Sprintf("%+v", subPayload), fmt.Sprintf("%+v", h))
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
		// TODO: send a chat message
	case SetHp:
		//   |<type>|<position spec>|<HP number>
		// 0 |  1   |       2       |     3
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		h := *new(HPStatus)
		hp, _ := strconv.Atoi(msg.Get(3))
		h.Current = hp
		delta := PokeDelta{HP: h}
		subPayload.PlayerId = p.PlayerId
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint("hp set room state payload", fmt.Sprintf("%+v", payload))
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case StatusInflict, StatusCure:
		//   |<type>|<position spec>|<status>
		// 0 |  1   |       2       |   3
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		h := HPStatus{Status: msg.Get(3)}
		delta := PokeDelta{HP: h}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint("major status update", msgType, h.Status)
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case TeamCure, InvertBoost, ClearBoost, ClearNegBoost, AbilityEnd:
		//   |<type>|<position spec>
		// 0 |  1   |       2
		// position spec is included to describe the user/cause of the team cure
		// for other messages, it is the target of the stat boost change
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint(msgType, "issued for", payload.Player.PlayerId)
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Boost, Unboost, SetBoost:
		//   |<type>|<position spec>|<stat spec>|<amount>
		// 0 |  1   |       2       |     3     |   4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		amount, _ := strconv.Atoi(msg.Get(4))
		b := StatMod{Stat: msg.Get(3), Amount: amount}
		delta := PokeDelta{Boost: b}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint(p, b.Stat, "boosted by", b.Amount)
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case SwapBoost:
		//   |<type>|<source position spec>|<target position spec>|<stat spec list>
		// 0 |  1   |           2          |           3          |       4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		t := NewPokemonPosition(msg.Get(3))
		s := msg.Get(4)
		b := StatMod{Stat: s, Reference: t}
		delta := PokeDelta{Boost: b}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case CopyBoost:
		//   |<type>|<source position spec>|<target position spec>
		// 0 |  1   |           2          |           3
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		t := NewPokemonPosition(msg.Get(3))
		b := StatMod{Reference: t}
		delta := PokeDelta{Boost: b}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case ClearPosBoost:
		//   |<type>|<target position spec>|<cause position spec>|<desc>
		// 0 |  1   |           2          |           3         |  4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case EffectStart, EffectEnd:
		//   |<type>|<position spec>|<effect>
		// 0 |  1   |       2       |   3
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		e := msg.Get(3)
		delta := PokeDelta{Effect: e}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		goPrint("effect", e, msgType, "on pokemon", p)
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Item, ItemEnd:
		//   |<type>|<position spec>|<item>|<optional: [from]>|<optional: [silent] (ItemEnd only)>
		// 0 |  1   |       2       |  3   |        4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		i := msg.Get(3)
		delta := PokeDelta{Item: i}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Ability:
		//   |<type>|<position spec>|<ability>|<optional [from]>
		// 0 |  1   |       2       |    3    |       4
		subPayload := new(UpdatePlayerPayload)
		p := NewPokemonPosition(msg.Get(2))
		ab := msg.Get(3)
		delta := PokeDelta{Ability: ab}
		subPayload.ActivePokemon = UpdatePlayerPokemon{Reason: msgType, Position: p, Delta: delta}
		payload := UpdateRoomStatePayload{RoomId: roomId, Player: *subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Transform:
		//
	case MegaEvolve:
		//
	case PrimalForm:
		//
	case Burst:
		//
	case SingleMoveEffect, SingleTurnEffect:
		//
	default:
		goPrint("not a supported minor battle action", msg.ReassembleTail(0))
	}
}

func (a *App) parseFieldBattleMessage(roomId string, msg *SplitString) {
	msgType := MessageType(msg.Get(1))
	switch msgType {
	case ClearAllBoost, SwapSideConditions:
		//   |<type>
		// 0 |  1
		subPayload := UpdateFieldPayload{Reason: msgType}
		payload := UpdateRoomStatePayload{RoomId: roomId, Field: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case Weather, FieldStart, FieldEnd:
		//   |<type>|<condition>
		// 0 |  1   |     2
		subPayload := UpdateFieldPayload{Reason: msgType, Condition: msg.Get(2)}
		payload := UpdateRoomStatePayload{RoomId: roomId, Field: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	case SideStart, SideEnd:
		//   |<type>|<side>|<condition>
		// 0 |  1   |  2   |     3
		goPrint("side start or side end message")
		goPrint("side is provided in this format:", msg.Get(2))
		subPayload := UpdateFieldPayload{Reason: msgType, PlayerId: msg.Get(2), Condition: msg.Get(3)}
		payload := UpdateRoomStatePayload{RoomId: roomId, Field: subPayload}
		a.channels.frontendChan <- ShowdownEvent{RoomStateTopic, payload}
	}
}

// produce a chat message only
// everything else lol
