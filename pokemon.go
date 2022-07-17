package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// for dealing with data representing individual and teams of pokemon

// a pokemon in json format
type UnpackedPokemon struct {
	Name      string         `json:"name"`      // optional
	Species   string         `json:"species"`   // required
	Gender    string         `json:"gender"`    // "M", "F" or blank
	Item      string         `json:"item"`      // optional
	Ability   string         `json:"ability"`   // could be "0" "1" or "H" or name of ability
	Evs       map[string]int `json:"evs"`       // each default to 0
	Ivs       map[string]int `json:"ivs"`       // each default to 31
	Nature    string         `json:"nature"`    // defaults to "Serious"
	Moves     []string       `json:"moves"`     // should be 1 to 4 moves
	Happiness int            `json:"happiness"` // leave blank for 255
	Shiny     string         `json:"shiny"`     // store value as "Yes" or blank
	Level     int            `json:"level"`     // leave blank for 100
	Pokeball  string         `json:"pokeball"`  // not even used?
}

type UnpackedPokemonList struct {
	Team []UnpackedPokemon
}

var statKeys = [6]string{"HP", "Atk", "Def", "SpA", "SpD", "Spe"}

func defaultPokemon() *UnpackedPokemon {
	return &UnpackedPokemon{Level: 100, Happiness: 255, Nature: "Serious", Moves: make([]string, 0)}
}

func parsePokeJson(poke []byte) *UnpackedPokemon {
	p := defaultPokemon()
	err := json.Unmarshal(poke, p)
	if err != nil {
		goPrint("could not parse poke json", string(poke))
		return nil
	}
	return p
}

func (p *UnpackedPokemon) pack() string {
	// represents a pokemon
	// NICKNAME|SPECIES|ITEM|ABILITY|MOVES|NATURE|EVS|GENDER|IVS|SHINY|LEVEL|HAPPINESS,POKEBALL,HIDDENPOWERTYPE
	// list delimiter is "]". no newlines.
	var builder strings.Builder
	// name and species
	if p.Name != "" {
		builder.WriteString(fmt.Sprintf("|%s|%s", p.Name, p.Species))
	} else {
		builder.WriteString(fmt.Sprintf("|%s|", p.Species))
	}
	// item and ability
	builder.WriteString(fmt.Sprintf("|%s", stringToPackedStyle(p.Item)))
	builder.WriteString(fmt.Sprintf("|%s", stringToPackedStyle(p.Ability)))
	// moves
	builder.WriteByte('|')
	for i, move := range p.Moves {
		if i == len(p.Moves)-1 {
			builder.WriteString(stringToPackedStyle(move))
		} else {
			builder.WriteString(fmt.Sprintf("%s,", stringToPackedStyle(move)))
		}
	}
	builder.WriteString(fmt.Sprintf("|%s", p.Nature))
	// evs
	builder.WriteByte('|')
	if len(p.Evs) > 0 {
		for i, stat := range statKeys {
			if i == 5 {
				builder.WriteString(fmt.Sprintf("%d,", p.Evs[stat]))
			} else {
				builder.WriteString(fmt.Sprintf("%d", p.Evs[stat]))
			}
		}
	}
	builder.WriteString(fmt.Sprintf("|%s", p.Gender))
	// ivs
	builder.WriteByte('|')
	if len(p.Ivs) > 0 {
		for i, stat := range statKeys {
			if i == 5 {
				builder.WriteString(fmt.Sprintf("%d,", p.Ivs[stat]))
			} else {
				builder.WriteString(fmt.Sprintf("%d", p.Ivs[stat]))
			}
		}
	}
	// shiny
	if strings.ToLower(p.Shiny) == "yes" {
		builder.WriteString("|S")
	} else {
		builder.WriteByte('|')
	}
	// level
	if p.Level != 100 {
		builder.WriteString(fmt.Sprintf("|%d", p.Level))
	} else {
		builder.WriteByte('|')
	}
	// happiness, pokeball, hp type
	builder.WriteByte('|')
	if p.Happiness != 255 || p.Pokeball != "" {
		builder.WriteString(fmt.Sprintf("%d,%s,", p.Happiness, p.Pokeball)) // TODO handle hidden power type
	}

	return builder.String()
}

func (p *UnpackedPokemon) toPokepaste() string {
	/*
		mostly same rules as the packed format
		Nickname <or species of no nickname> (species <if nickname>) (gender <optional>) @ item
		Ability: ability
		Level: level <line absent if 100>
		Shiny: Yes <or line absent>
		EVs: </ delimited with number first, then 3 letter abbreviation. 0s absent>
		nature Nature
		IVs: <same as Evs>
		- <each>
		- <move>
		- <listed>
		- <like this>

	*/
	var builder strings.Builder

	return builder.String()
}

func parsePokepaste(paste string) *UnpackedPokemon {
	p := defaultPokemon()
	// parse the pokepaste
	pasteLines := strings.Split(paste, "\n")
	// first line

	var leftovers string
	itemSplit := NewSplitString(pasteLines[0], " @ ")
	if itemSplit.len >= 2 {
		if strings.Contains(itemSplit.Get(-1), "(") {
			// @ signs are in the name only and there is no item
			// example lewis@ (Growlithe) (M)
			leftovers = pasteLines[0]
		} else {
			p.Item = strings.TrimSpace(itemSplit.Get(-1))
			leftovers = itemSplit.ReassembleHead(-1)
		}
	} else {
		leftovers = pasteLines[0]
	}
	// split on open parens
	// check how many resulting groups
	// last one may be gender
	// second to last must be species, unless last is
	// remained are part of nickname (what kind of sick fuck puts parens in their pokemon nicknames)
	// reassemble nickname
	nameSpeciesGender := strings.TrimSpace(leftovers)
	nsgParts := NewSplitString(nameSpeciesGender, "(")
	if nsgParts.len >= 3 {
		// nickname, species and maybe gender
		if len(nsgParts.Get(-1)) == 2 {
			p.Name = nsgParts.ReassembleHead(-2)
			p.Species = strings.Replace(nsgParts.Get(-2), ")", "", 1) // second to last
			p.Gender = strings.Replace(nsgParts.Get(-1), ")", "", 1)  // last group
		} else {
			p.Name = nsgParts.ReassembleHead(-1)
			p.Species = strings.Replace(nsgParts.Get(-1), ")", "", 1) // second to last
		}

	} else if nsgParts.len == 2 {
		// nickname and species or species and gender
		if len(nsgParts.Get(1)) == 2 {
			p.Species = nsgParts.Get(0)
			p.Gender = strings.Replace(nsgParts.Get(1), ")", "", 1)
		} else {
			p.Name = nsgParts.Get(0)
			p.Species = strings.Replace(nsgParts.Get(1), ")", "", 1)
		}
	} else {
		p.Species = nameSpeciesGender
	}
	// iterate of remaining lines
	for _, line := range pasteLines[1:] {
		chunks := NewSplitString(line, " ")
		switch chunks.Get(0) {
		case "-":
			p.Moves = append(p.Moves, chunks.ReassembleTail(1))
		case "Ability":
			p.Ability = chunks.ReassembleTail(1)
		case "Level":
			p.Level, _ = strconv.Atoi(chunks.Get(1))
		case "Happiness":
			p.Happiness, _ = strconv.Atoi(chunks.Get(1))
		case "EVs":
			p.Evs = parsePokepasteStatSpread(make(map[string]int), chunks.ReassembleTail(1))
		case "IVs":
			p.Ivs = parsePokepasteStatSpread(defaultIvs(), chunks.ReassembleTail(1))
		case "Shiny":
			p.Shiny = chunks.Get(1) // should say "Yes", any other value will be ignored when generating other formats
		case "Pokeball":
			p.Pokeball = chunks.ReassembleTail(1)
		default:
			if chunks.Get(1) == "Nature" {
				p.Nature = chunks.Get(0)
			}
		}
	}
	return p
}

func stringToPackedStyle(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}

func defaultIvs() map[string]int {
	ivs := make(map[string]int)
	for _, k := range statKeys {
		ivs[k] = 31
	}
	return ivs
}

//
func parsePokepasteStatSpread(defaults map[string]int, stats string) map[string]int {
	splitStats := NewSplitString(stats, "/")
	for _, s := range splitStats.inner {
		valueAndLabel := NewSplitString(strings.TrimSpace(s), " ")
		for _, key := range statKeys {
			if strings.ToLower(key) == strings.TrimSpace(strings.ToLower(valueAndLabel.Get(1))) {
				defaults[key], _ = strconv.Atoi(valueAndLabel.Get(0))
			}
		}
	}

	return defaults
}
