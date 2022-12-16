package main

import (
	"encoding/json"

	"golang.org/x/exp/slices"
)

func reconcileRoomStateInner(update UpdateRoomStatePayload, state RoomState) RoomState {
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
			state.Participants[update.Player.PlayerId] = reconcileFrontendPlayerState(update.Player, state.Participants[update.Player.PlayerId])
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
	return state
}

func reconcileFrontendPlayerState(update UpdatePlayerPayload, state BattleRoomParticipant) BattleRoomParticipant {
	if update.TeamSize != 0 {
		state.TeamSize = update.TeamSize
	}
	if update.ActivePokemon.Position.NickName != "" {
		switch update.ActivePokemon.Reason {
		case Switch, Drag:
			position := update.ActivePokemon.Position.Position
			for i := len(state.Active); i < position+1; i++ {
				state.Active = append(state.Active, PokemonState{})
			}
			activeSwitchingOut := state.Active[position]
			activeSwitchingOut.Active = false
			switchingIn := pokeStateFromUpdate(update.ActivePokemon)
			if activeSwitchingOut.Species == "" {
				// this means the pokemon is switching into a new position
				state.Active[position] = switchingIn
			} else {
				// find pokemon that matches switching out on inactive list
				switchIndex := findPoke(switchingIn, state.Inactive)
				if switchIndex >= 0 {
					// swap the two pokemon
					switchingIn = state.Inactive[switchIndex]
					switchingIn.Active = true
					state.Inactive[switchIndex] = activeSwitchingOut
					state.Active[position] = switchingIn
				} else {
					// if not on inactive list, then just place it there and create new state for switch in
					state.Inactive = append(state.Inactive, activeSwitchingOut)
					state.Active[position] = switchingIn
				}
			}
		case Damage, Heal:
			position := update.ActivePokemon.Position.Position
			active := state.Active[position]
			active.CurrentHp = update.ActivePokemon.Delta.HP.Current
			active.MaxHp = update.ActivePokemon.Delta.HP.Max
			active.MajorStatus = update.ActivePokemon.Delta.HP.Status
			state.Active[position] = active
		case SetHp:
			position := update.ActivePokemon.Position.Position
			active := state.Active[position]
			active.CurrentHp = update.ActivePokemon.Delta.HP.Current
			state.Active[position] = active
		case Faint:
			position := update.ActivePokemon.Position.Position
			state.Active[position].Fainted = true
		case StatusInflict, StatusCure:
			position := update.ActivePokemon.Position.Position
			state.Active[position].MajorStatus = update.ActivePokemon.Delta.HP.Status
		case TeamCure:
			for _, ap := range state.Active {
				ap.MajorStatus = ""
			}
			for _, ip := range state.Inactive {
				ip.MajorStatus = ""
			}
		case Boost, Unboost:
			statKey := update.ActivePokemon.Delta.Boost.Stat
			statChange := update.ActivePokemon.Delta.Boost.Amount
			position := update.ActivePokemon.Position.Position
			stats := state.Active[position].StatBoosts
			currentBoost, boostPresent := stats[statKey]
			if !boostPresent {
				stats[statKey] = statChange
			} else {
				stats[statKey] = currentBoost + statChange
			}
			state.Active[position].StatBoosts = stats
		case SetBoost:
			statKey := update.ActivePokemon.Delta.Boost.Stat
			statChange := update.ActivePokemon.Delta.Boost.Amount
			position := update.ActivePokemon.Position.Position
			state.Active[position].StatBoosts[statKey] = statChange
		case InvertBoost:
			position := update.ActivePokemon.Position.Position
			stats := state.Active[position].StatBoosts
			for stat, boost := range stats {
				state.Active[position].StatBoosts[stat] = -boost
			}
		case ClearBoost:
			position := update.ActivePokemon.Position.Position
			stats := state.Active[position].StatBoosts
			for stat := range stats {
				state.Active[position].StatBoosts[stat] = 0
			}
		case ClearPosBoost:
			position := update.ActivePokemon.Position.Position
			stats := state.Active[position].StatBoosts
			for stat, boost := range stats {
				if boost > 0 {
					state.Active[position].StatBoosts[stat] = 0
				}
			}
		case ClearNegBoost:
			position := update.ActivePokemon.Position.Position
			stats := state.Active[position].StatBoosts
			for stat, boost := range stats {
				if boost < 0 {
					state.Active[position].StatBoosts[stat] = 0
				}
			}
		case EffectStart:
			effect := update.ActivePokemon.Delta.Effect
			position := update.ActivePokemon.Position.Position
			state.Active[position].MinorStatuses = append(state.Active[position].MinorStatuses, effect)
		case EffectEnd:
			effect := update.ActivePokemon.Delta.Effect
			position := update.ActivePokemon.Position.Position
			effects := state.Active[position].MinorStatuses
			effectInd := slices.Index(effects, effect)
			state.Active[position].MinorStatuses = slices.Delete(effects, effectInd, effectInd+1)
		case Item:
			position := update.ActivePokemon.Position.Position
			item := update.ActivePokemon.Delta.Item
			state.Active[position].HeldItem = item
		case ItemEnd:
			position := update.ActivePokemon.Position.Position
			state.Active[position].HeldItem = ""
		case Ability:
			position := update.ActivePokemon.Position.Position
			ability := update.ActivePokemon.Delta.Ability
			state.Active[position].Ability = ability
		case AbilityEnd:
			position := update.ActivePokemon.Position.Position
			state.Active[position].Ability = ""
		}
	}
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
