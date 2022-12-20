package backend

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

func reconcileRoomStateInner(update UpdateRoomStatePayload, state RoomState) RoomState {
	reason := update.Player.ActivePokemon.Reason
	if string(reason) == "" {
		reason = update.Field.Reason
	}
	goPrint("reconciling room state for reason: ", reason)
	goPrint("[ UPDATE", reason, "]", fmt.Sprintf("%+v", updateWithoutRequest(update)))
	goPrint("[ BASE", reason, "]", fmt.Sprintf("%+v", roomStateWithoutRequest(state)))
	if update.Title != "" {
		state.Title = update.Title
	}
	if update.Gen != 0 {
		state.Gen = update.Gen
	}
	if update.GameType != "" {
		state.GameType = update.GameType
	}
	if update.Request != "" {
		requestJson := make(map[string]interface{})
		err := json.Unmarshal([]byte(update.Request), &requestJson)
		if err != nil {
			goPrint("error unmarshalling request", update.Request, err.Error())
		}
		state.Request = requestJson
	}
	if update.Tier != "" {
		state.Tier = update.Tier
	}
	if update.Rated {
		state.Rated = update.Rated
	}
	if update.Active != nil {
		state.Active = *update.Active
	}
	if update.Timer != nil {
		state.Timer = *update.Timer
	}
	if update.Player.PlayerId != "" {
		if state.Participants == nil {
			state.Participants = make(map[string]BattleRoomParticipant)
		}
		if _, hasKey := state.Participants[update.Player.PlayerId]; hasKey {
			state.Participants = reconcileFrontendPlayerState(update.Player, state.Participants)
		} else {
			newPlayer := new(BattleRoomParticipant)
			newPlayer.Name = update.Player.Name
			newPlayer.Avatar = update.Player.Avatar
			newPlayer.Rating = update.Player.Rating
			newPlayer.Active = make([]PokemonState, 0)
			newPlayer.Inactive = make([]PokemonState, 0)
			newPlayer.Id = update.Player.PlayerId
			state.Participants[update.Player.PlayerId] = *newPlayer
		}
	}
	if update.Field.Reason != "" {
		if state.Field.Sides == nil {
			state.Field.Sides = make(map[string][]BattleFieldCondition)
		}
		switch update.Field.Reason {
		case Weather:
			weatherIndex := slices.IndexFunc(state.Field.Conditions, func(bfc BattleFieldCondition) bool { return bfc.Weather })
			if update.Field.Condition == "none" {
				state.Field.Conditions = slices.Delete(state.Field.Conditions, weatherIndex, weatherIndex+1)
			} else if weatherIndex != -1 {
				state.Field.Conditions[weatherIndex].Condition = update.Field.Condition
			} else {
				state.Field.Conditions = append(state.Field.Conditions, BattleFieldCondition{Condition: update.Field.Condition, Weather: true})
			}
		case FieldStart:
			state.Field.Conditions = append(state.Field.Conditions, BattleFieldCondition{Condition: update.Field.Condition, Weather: false})
		case FieldEnd:
			conditionIndex := slices.IndexFunc(state.Field.Conditions, func(bfc BattleFieldCondition) bool { return update.Field.Condition == bfc.Condition })
			state.Field.Conditions = slices.Delete(state.Field.Conditions, conditionIndex, conditionIndex+1)
		case ClearAllBoost:
			for _, player := range state.Participants {
				for _, pokemon := range player.Active {
					pokemon.StatBoosts = make(map[string]int)
				}
			}
		case SideStart:
			_, hasKey := state.Field.Sides[update.Field.PlayerId]
			c := BattleFieldCondition{Condition: update.Field.Condition, Weather: false}
			if hasKey {
				state.Field.Sides[update.Field.PlayerId] = append(state.Field.Sides[update.Field.PlayerId], c)
			} else {
				state.Field.Sides[update.Field.PlayerId] = []BattleFieldCondition{c}
			}
		case SideEnd:
			conditionIndex := slices.IndexFunc(state.Field.Sides[update.Field.PlayerId], func(bfc BattleFieldCondition) bool { return update.Field.Condition == bfc.Condition })
			state.Field.Sides[update.Field.PlayerId] = slices.Delete(state.Field.Conditions, conditionIndex, conditionIndex+1)
		}
	}
	goPrint("[ RESULT", reason, "]", fmt.Sprintf("%+v", roomStateWithoutRequest(state)))
	return state
}

func reconcileFrontendPlayerState(update UpdatePlayerPayload, state map[string]BattleRoomParticipant) map[string]BattleRoomParticipant {
	playerState := state[update.PlayerId]
	if update.TeamSize != 0 {
		playerState.TeamSize = update.TeamSize
	}
	if update.ActivePokemon.Position.NickName != "" {
		switch update.ActivePokemon.Reason {
		case Switch, Drag:
			position := update.ActivePokemon.Position.Position
			for i := len(playerState.Active); i < position+1; i++ {
				playerState.Active = append(playerState.Active, PokemonState{})
			}
			activeSwitchingOut := playerState.Active[position]
			activeSwitchingOut.Active = false
			activeSwitchingOut.StatBoosts = make(map[string]int)
			activeSwitchingOut.MinorStatuses = make([]string, 0)
			switchingIn := pokeStateFromUpdate(update.ActivePokemon)
			if activeSwitchingOut.Species == "" {
				// this means the pokemon is switching into a new position
				playerState.Active[position] = switchingIn
			} else {
				// find pokemon that matches switching out on inactive list
				switchIndex := findPoke(switchingIn, playerState.Inactive)
				if switchIndex >= 0 {
					// swap the two pokemon
					switchingIn = playerState.Inactive[switchIndex]
					switchingIn.Active = true
					playerState.Inactive[switchIndex] = activeSwitchingOut
					playerState.Active[position] = switchingIn
				} else {
					// if not on inactive list, then just place it there and create new playerState for switch in
					playerState.Inactive = append(playerState.Inactive, activeSwitchingOut)
					playerState.Active[position] = switchingIn
				}
			}
		case Damage, Heal:
			position := update.ActivePokemon.Position.Position
			active := playerState.Active[position]
			active.CurrentHp = update.ActivePokemon.Delta.HP.Current
			active.MaxHp = update.ActivePokemon.Delta.HP.Max
			active.MajorStatus = update.ActivePokemon.Delta.HP.Status
			playerState.Active[position] = active
		case SetHp:
			position := update.ActivePokemon.Position.Position
			active := playerState.Active[position]
			active.CurrentHp = update.ActivePokemon.Delta.HP.Current
			playerState.Active[position] = active
		case Faint:
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].Fainted = true
		case StatusInflict, StatusCure:
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].MajorStatus = update.ActivePokemon.Delta.HP.Status
		case TeamCure:
			for _, ap := range playerState.Active {
				ap.MajorStatus = ""
			}
			for _, ip := range playerState.Inactive {
				ip.MajorStatus = ""
			}
		case Boost, Unboost:
			statKey := update.ActivePokemon.Delta.Boost.Stat
			boostSign := 1
			if update.ActivePokemon.Reason == Unboost {
				boostSign = -1
			}
			statChange := update.ActivePokemon.Delta.Boost.Amount
			position := update.ActivePokemon.Position.Position
			stats := playerState.Active[position].StatBoosts
			currentBoost, boostPresent := stats[statKey]
			goPrint("boosting", playerState.Id, position, statKey, statChange, currentBoost)
			if !boostPresent {
				stats[statKey] = (statChange * boostSign)
			} else {
				stats[statKey] = currentBoost + (statChange * boostSign)
			}
			playerState.Active[position].StatBoosts = stats
		case SetBoost:
			statKey := update.ActivePokemon.Delta.Boost.Stat
			statChange := update.ActivePokemon.Delta.Boost.Amount
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].StatBoosts[statKey] = statChange
		case InvertBoost:
			position := update.ActivePokemon.Position.Position
			stats := playerState.Active[position].StatBoosts
			for stat, boost := range stats {
				playerState.Active[position].StatBoosts[stat] = -boost
			}
		case ClearBoost:
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].StatBoosts = make(map[string]int)
		case ClearPosBoost:
			position := update.ActivePokemon.Position.Position
			stats := playerState.Active[position].StatBoosts
			for stat, boost := range stats {
				if boost > 0 {
					playerState.Active[position].StatBoosts[stat] = 0
				}
			}
		case ClearNegBoost:
			position := update.ActivePokemon.Position.Position
			stats := playerState.Active[position].StatBoosts
			for stat, boost := range stats {
				if boost < 0 {
					playerState.Active[position].StatBoosts[stat] = 0
				}
			}
		case CopyBoost:
			targetPosition := update.ActivePokemon.Delta.Boost.Reference.Position
			targetOwnerId := update.ActivePokemon.Delta.Boost.Reference.PlayerId
			targetOwner := state[targetOwnerId]
			targetBoosts := targetOwner.Active[targetPosition].StatBoosts
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].StatBoosts = targetBoosts
		case SwapBoost:
			statList := strings.Split(update.ActivePokemon.Delta.Boost.Stat, ",")
			targetPosition := update.ActivePokemon.Delta.Boost.Reference.Position
			targetOwnerId := update.ActivePokemon.Delta.Boost.Reference.PlayerId
			targetOwner := state[targetOwnerId]
			target := targetOwner.Active[targetPosition]
			position := update.ActivePokemon.Position.Position
			active := playerState.Active[position]
			for _, stat := range statList {
				tempBoost := target.StatBoosts[stat]
				target.StatBoosts[stat] = active.StatBoosts[stat]
				active.StatBoosts[stat] = tempBoost
			}
			state[targetOwnerId].Active[targetPosition] = target
			playerState.Active[position] = active
		case EffectStart:
			effect := update.ActivePokemon.Delta.Effect
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].MinorStatuses = append(playerState.Active[position].MinorStatuses, effect)
		case EffectEnd:
			effect := update.ActivePokemon.Delta.Effect
			position := update.ActivePokemon.Position.Position
			effects := playerState.Active[position].MinorStatuses
			effectInd := slices.Index(effects, effect)
			playerState.Active[position].MinorStatuses = slices.Delete(effects, effectInd, effectInd+1)
		case Item:
			position := update.ActivePokemon.Position.Position
			item := update.ActivePokemon.Delta.Item
			playerState.Active[position].HeldItem = item
		case ItemEnd:
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].HeldItem = ""
		case Ability:
			position := update.ActivePokemon.Position.Position
			ability := update.ActivePokemon.Delta.Ability
			playerState.Active[position].Ability = ability
		case AbilityEnd:
			position := update.ActivePokemon.Position.Position
			playerState.Active[position].Ability = ""
		}
	}
	state[update.PlayerId] = playerState
	return state
}

func pokeStateFromUpdate(p UpdatePlayerPokemon) PokemonState {
	s := new(PokemonState)
	s.Species = p.Details.Species
	s.Level = p.Details.Level
	s.Gender = p.Details.Gender
	s.Shiny = p.Details.Shiny
	s.NickName = p.Position.NickName
	s.PlayerId = p.Position.PlayerId
	s.CurrentHp = p.Delta.HP.Current
	s.MaxHp = p.Delta.HP.Max
	s.MajorStatus = p.Delta.HP.Status
	s.Active = true
	s.Fainted = false
	s.StatBoosts = make(map[string]int)
	s.Moves = make([]interface{}, 0)
	s.MinorStatuses = make([]string, 0)
	return *s
}

// check if two states are equal
func equalPokeState(a, b PokemonState) bool {
	if a.Species == "" || b.Species == "" {
		return false
	}
	return (a.Level == b.Level &&
		a.Gender == b.Gender &&
		a.Shiny == b.Shiny &&
		a.Species == b.Species &&
		a.NickName == b.NickName)
}

// look for matching pokemon in slice and return its index (-1 if not found)
func findPoke(p PokemonState, s []PokemonState) int {
	for i, b := range s {
		if equalPokeState(p, b) {
			return i
		}
	}
	return -1
}

func roomStateWithoutRequest(r RoomState) RoomState {
	r.Request = nil
	return r
}

func updateWithoutRequest(u UpdateRoomStatePayload) UpdateRoomStatePayload {
	u.Request = ""
	return u
}
