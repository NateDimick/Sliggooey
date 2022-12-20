package backend

import "testing"

// Testing strategy for roomstate:
// create json resources of room state and room update documents
// unmarshal them and pass them in to reconcileRoomState
// test against expected outputs

// base state is a starter room state to update and check against
// it can grow as unit tests demand more comparisons
var baseState RoomState = RoomState{
	Participants: map[string]BattleRoomParticipant{
		"p1a": {
			Id: "p1a",
			Active: []PokemonState{
				{
					Species:    "pikachu",
					NickName:   "unit test",
					StatBoosts: make(map[string]int),
				},
			},
		},
		"p2b": {
			Id: "p2b",
		},
	},
	Field: BattleFieldState{
		Conditions: make([]BattleFieldCondition, 0),
		Sides:      make(map[string][]BattleFieldCondition),
	},
}

func TestBoostFromNeutral(t *testing.T) {
	update := UpdateRoomStatePayload{
		Player: UpdatePlayerPayload{
			PlayerId: "p1a",
			ActivePokemon: UpdatePlayerPokemon{
				Reason: Boost,
				Position: PokemonPosition{
					PlayerId: "p1a",
					Position: 0,
					NickName: "unit test",
				},
				Delta: PokeDelta{
					Boost: StatMod{
						Stat:   "atk",
						Amount: 1,
					},
				},
			},
		},
	}
	result := reconcileRoomStateInner(update, baseState)
	valueAfterBoost := result.Participants["p1a"].Active[0].StatBoosts["atk"]
	if valueAfterBoost != 1 {
		t.Log("stat did not boost properly!")
		t.Fail()
	}
}

func TestBoostBoostedStatFurther(t *testing.T) {

}

func TestUnboostFromNeutral(t *testing.T) {

}

func TestUnboostBoostedStat(t *testing.T) {

}

func TestAddWeather(t *testing.T) {
	update := UpdateRoomStatePayload{
		Field: UpdateFieldPayload{
			Condition: "rain",
			Reason:    Weather,
		},
	}
	result := reconcileRoomStateInner(update, baseState)
	fieldAfterUpdate := result.Field
	if len(fieldAfterUpdate.Conditions) != 1 {
		t.Log("condition not added to field list")
		t.Fail()
	}
	if fieldAfterUpdate.Conditions[0].Condition != "rain" {
		t.Log("condition value not set to rain")
		t.Fail()
	}
	if !fieldAfterUpdate.Conditions[0].Weather {
		t.Log("field condition not marked as weather")
		t.Fail()
	}
}

func TestChangeWeather(t *testing.T) {

}

func TestEndWeather(t *testing.T) {

}

func TestAddFieldCondition(t *testing.T) {

}

func TestAddAnotherFieldCondition(t *testing.T) {

}

func TestEndFieldCondition(t *testing.T) {

}

func TestAddSideCondition(t *testing.T) {
	update := UpdateRoomStatePayload{
		Field: UpdateFieldPayload{
			Reason:    SideStart,
			Condition: "spikes",
			PlayerId:  "p1a",
		},
	}
	result := reconcileRoomStateInner(update, baseState)
	if _, ok := result.Field.Sides["p1a"]; !ok {
		t.Log("side not added to sides map")
		t.Fail()
	}
	if len(result.Field.Sides["p1a"]) != 1 {
		t.Log("condition not added to array")
		t.Fail()
	}
	if result.Field.Sides["p1a"][0].Condition != "spikes" {
		t.Log("condition not filled correctly")
		t.Fail()
	}

}

func TestAddAnotherSideCondition(t *testing.T) {

}

func TestEndSideCondition(t *testing.T) {

}
