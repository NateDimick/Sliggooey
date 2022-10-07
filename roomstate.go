package main

import (
	"encoding/json"
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
			active.CurrentHp = update.ActivePokemon.HP.Current
			active.MaxHp = update.ActivePokemon.HP.Max
			active.MajorStatus = update.ActivePokemon.HP.Status
			state.Active[position] = active
		case Faint:
			position := update.ActivePokemon.Position.Position
			state.Active[position].Fainted = true
		case Boost, Unboost:
			//
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
	s.CurrentHp = p.HP.Current
	s.MaxHp = p.HP.Max
	s.MajorStatus = p.HP.Status
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