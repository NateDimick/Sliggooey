package main

import (
	"encoding/json"
	"testing"
)

// this is a bad test. re-write.
func TestParseSearchJson(t *testing.T) {
	cases := []struct {
		input  string
		expect GamesStatus
	}{
		{"{\"searching\":[],\"games\":null}", GamesStatus{}},
		{"{\"searching\":[],\"games\":{\"battle-gen8randombattle-1618577741\":\"[Gen 8] Random Battle\"}}", GamesStatus{}},
	}
	for _, c := range cases {
		gs := new(GamesStatus)
		err := json.Unmarshal([]byte(c.input), gs)
		if err != nil {
			t.Fatalf("got an error: %s", err.Error())
		} else {
			t.Logf("%s", gs)
			t.Logf("%s", gs.Games)
			t.Fail()
		}
	}
}

func TestParseRequestJson(t *testing.T) {

}
