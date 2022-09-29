package main

import (
	"testing"
)

func TestMakeBattleChoice(t *testing.T) {
	tests := []struct {
		choice BattleChoice
		expect string
	}{
		{BattleChoice{Attack, 2, "", 1, "mega"}, "/choose move 2 +1 mega"},
		{BattleChoice{Attack, 1, "", 0, ""}, "/choose move 1"},
		{BattleChoice{SwitchOut, 0, "", 4, ""}, "/choose switch 4"},
		{BattleChoice{Default, 0, "", 0, ""}, "/choose default"},
		{BattleChoice{Pass, 0, "", 0, ""}, "/choose pass"},
		{BattleChoice{Undo, 0, "", 0, ""}, "/choose undo"},
	}

	for _, testcase := range tests {
		result := FormatBattleChoices(testcase.choice)
		if result != testcase.expect {
			t.Fatalf("expected chose command '%s' but got '%s'", testcase.expect, result)
		}
	}

}

func TestMakeMultiBattleChoice(t *testing.T) {
	a := BattleChoice{Attack, 1, "", 0, ""}
	b := BattleChoice{SwitchOut, 0, "", 4, ""}
	expect := "/choose move 1, switch 4"
	result := FormatBattleChoices(a, b)
	if result != expect {
		t.Fatalf("expected chose command '%s' but got '%s'", expect, result)

	}

}

func TestBuildCommand(t *testing.T) {
	expect := "/challenge me, battle"
	result := buildCommand(Challenge, "me", "battle")
	if result != expect {
		t.Fatalf("expected %s, got %s", expect, result)
	}
}
