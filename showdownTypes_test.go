package main

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, desc string, expected, actual interface{}) {
	if expected != actual {
		t.Fatalf(fmt.Sprintf("FAILED: %s: expected %v but got %v", desc, expected, actual))
	}
}

func TestNewPokemonPosition(t *testing.T) {
	tests := []struct {
		input  string
		output *PokemonPosition
	}{
		{"p1a: Dragonite", &PokemonPosition{"p1", 1, "Dragonite"}},
		{"p2b: Garbodor", &PokemonPosition{"p2", 2, "Garbodor"}},
		{"p3: Vanilluxe", &PokemonPosition{"p3", 0, "Vanilluxe"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := NewPokemonPosition(tt.input)
			assertEqual(t, "Player Id", tt.output.PlayerId, p.PlayerId)
			assertEqual(t, "Pokemon Position", tt.output.Position, p.Position)
			assertEqual(t, "Nickname", tt.output.NickName, p.NickName)
		})
	}

}

func TestNewPokemonDetails(t *testing.T) {
	tests := []struct {
		input  string
		output *PokemonDetails
	}{
		{"Salamence, L77, F", &PokemonDetails{"Salamence", 77, 'F', false}},
		{"Abomasnow, M, shiny", &PokemonDetails{"Abomasnow", 100, 'M', true}},
		{"Metagross", &PokemonDetails{"Metagross", 100, ' ', false}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			d := NewPokemonDetails(tt.input)
			assertEqual(t, "Species", tt.output.Species, d.Species)
			assertEqual(t, "level", tt.output.Level, d.Level)
			assertEqual(t, "gender", tt.output.Gender, d.Gender)
			assertEqual(t, "shiny", tt.output.Shiny, d.Shiny)
		})
	}
}
