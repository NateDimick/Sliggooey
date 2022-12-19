package backend

import (
	"fmt"
	"strings"
)

// types for the go backend

type AppError int

const (
	WebsocketError AppError = iota
	JsonUnmarshalError
)

type ShowdownEventTopic string

const (
	LoginFail         ShowdownEventTopic = "loginFail"
	LoginSuccess      ShowdownEventTopic = "loginSuccess"
	FormatTopic       ShowdownEventTopic = "formats"
	PMTopic           ShowdownEventTopic = "pm"
	NewRoomTopic      ShowdownEventTopic = "newRoom"
	RoomExitTopic     ShowdownEventTopic = "roomExit"
	RoomMessageTopic  ShowdownEventTopic = "roomMsg"
	RoomStateTopic    ShowdownEventTopic = "roomState"
	PopupTopic        ShowdownEventTopic = "popup"
	CurrentGamesTopic ShowdownEventTopic = "games"
	ChallengeTopic    ShowdownEventTopic = "challenged"
	ChallengeEndTopic ShowdownEventTopic = "challengeEnd"
)

type AppSettings struct {
}

type User struct {
	UserName string
	Rank     UserRank
	IsUser   bool
	IsAway   bool
	Status   string
}

func NewUser(nameStr string) *User {
	u := new(User)
	splitName := NewSplitString(nameStr[1:], "@")
	u.UserName = splitName.Get(0)
	u.Rank = UserRank(nameStr[0])
	u.IsUser = false
	if splitName.Get(1) != "" {
		u.IsAway = splitName.Get(1)[0] == '!'
	}
	if len(splitName.Get(1)) > 1 {
		u.Status = splitName.Get(1)[1:]
	}
	return u
}

type ShowdownCredentials struct {
	UserName string
	Password string
}

// an event that is sent to the front end by event emitter
type ShowdownEvent struct {
	Topic ShowdownEventTopic
	Data  interface{}
}

type LoginSuccessPayload struct {
	UserName string
}

type LoginFailurePayload struct {
	Reason string
}

type BattleChoice struct {
	Action  BattleAction
	Move    int    // could be a string, but I think this is cleaner. 0 is also invalid here
	AltMove string // only for special cases, like meta games with teams of more than 10 pokemon
	Target  int    // it's okay for target to default to 0 because 0 isn't a valid target. both for multi-battle target and switching
	Gimmick string
}

func (b *BattleChoice) format() string {
	switch b.Action {
	case TeamOrder:
		if b.Move > 0 {
			return fmt.Sprintf("%s %d", b.Action, b.Move)
		} else {
			return fmt.Sprintf("%s %s", b.Action, b.AltMove)
		}
	case SwitchOut:
		return fmt.Sprintf("%s %d", b.Action, b.Target)
	case Attack:
		if b.Target != 0 {
			targetMod := ""
			if b.Target > 0 {
				targetMod = "+"
			}
			return strings.TrimSpace(fmt.Sprintf("%s %d %s%d %s", b.Action, b.Move, targetMod, b.Target, b.Gimmick))
		} else {
			return strings.TrimSpace(fmt.Sprintf("%s %d %s", b.Action, b.Move, b.Gimmick))

		}
	case Default, Pass, Undo:
		return string(b.Action)
	default:
		return string(Default)
	}
}

type BattleAction string

const (
	TeamOrder BattleAction = "team"
	SwitchOut BattleAction = "switch"
	Attack    BattleAction = "move"
	Default   BattleAction = "default"
	Pass      BattleAction = "pass"
	Undo      BattleAction = "undo"
)

type PmSource int

const (
	Other PmSource = iota
	Self
	System
)

type PrivateMessagePayload struct {
	With    string
	IAm     string
	Message string
	From    PmSource
}

type NewRoomPayload struct {
	RoomId   string
	RoomType RoomType
}

type RoomMessagePayload struct {
	RoomId  string
	From    string
	Message string
}

type RoomHtmlPayload struct {
	RoomId string
	Html   string
	Name   string
	Update bool
}

type ChallengePayload struct {
	With       string
	IAm        string
	Format     string
	Challenger PmSource
}

type UpdateRoomStatePayload struct {
	RoomId   string              `json:"RoomId"`
	Title    string              `json:"title"`
	Gen      int                 `json:"gen"`
	GameType string              `json:"gameType"`
	Tier     string              `json:"tier"`
	Timer    *bool               `json:"timer"` // bools default to false, but pointers default to nil
	Rated    bool                `json:"rated"`
	Active   *bool               `json:"active"`
	Player   UpdatePlayerPayload `json:"player"`
	Request  string              `json:"request"`
	Field    UpdateFieldPayload  `json:"field"`
}

type UpdatePlayerPayload struct {
	PlayerId      string              `json:"playerId"`
	Name          string              `json:"playerName"`
	Avatar        string              `json:"avatarName"`
	Rating        string              `json:"rating"`
	TeamSize      int                 `json:"teamSize"`
	ActivePokemon UpdatePlayerPokemon `json:"activePokemonUpdate"`
}

type UpdatePlayerPokemon struct {
	Reason   MessageType     `json:"updateReason"`
	Position PokemonPosition `json:"positionalDetails"`
	Details  PokemonDetails  `json:"intrinsicDetails"`
	Delta    PokeDelta       `json:"delta"`
}

type UpdateFieldPayload struct {
	Reason    MessageType `json:"updateReason"`
	Condition string      `json:"condition"`
	PlayerId  string      `json:"side"` // present if the condition only effects one player
}

type SplitString struct {
	delimiter string
	inner     []string
	len       int
}

func NewSplitString(s string, delimiter string) *SplitString {
	ss := new(SplitString)
	ss.delimiter = delimiter
	ss.inner = strings.Split(s, delimiter)
	ss.len = len(ss.inner)
	return ss
}

// get string at given position. can be negative to index from rear. "" if index OOB
func (s *SplitString) Get(index int) string {
	if index < s.len && index >= 0 {
		return s.inner[index]
	} else if index < 0 && index > -s.len {
		return s.inner[s.len+index]
	} else {
		return ""
	}
}

// rejoin from a certain index, inclusive. can be negative
func (s *SplitString) ReassembleTail(from int) string {
	if from >= 0 {
		return strings.Join(s.inner[from:], s.delimiter)
	} else {
		return strings.Join(s.inner[s.len+from:], s.delimiter)
	}

}

// rejoin up to a certain index, exclusive. can be negative. "" if 0
func (s *SplitString) ReassembleHead(to int) string {
	if to > 0 {
		return strings.Join(s.inner[:to], s.delimiter)
	} else if to < 0 {
		return strings.Join(s.inner[:s.len+to], s.delimiter)
	} else {
		return ""
	}
}

// rejoin midsection. both must be positive
func (s *SplitString) ReassembleMid(from int, to int) string {
	if from < to {
		return strings.Join(s.inner[from:to], s.delimiter)
	} else {
		return ""
	}
}

type RoomState struct {
	Request      map[string]interface{}           `json:"request"`
	Title        string                           `json:"title"`
	Gen          int                              `json:"gen"`
	Tier         string                           `json:"tier"`
	GameType     string                           `json:"gameType"`
	Timer        bool                             `json:"timer"`
	Rated        bool                             `json:"rated"`
	Active       bool                             `json:"active"`
	Participants map[string]BattleRoomParticipant `json:"participants"`
	Field        BattleFieldState                 `json:"field"`
}

type BattleRoomParticipant struct {
	Name     string         `json:"name"`
	Id       string         `json:"id"`
	Rating   string         `json:"rating"`
	TeamSize int            `json:"teamSize"`
	Avatar   string         `json:"avatar"`
	Active   []PokemonState `json:"active"`
	Inactive []PokemonState `json:"inactive"`
}

type PokemonState struct {
	Species       string         `json:"species"`
	NickName      string         `json:"nickname"`
	Gender        rune           `json:"gender"`
	Level         int            `json:"level"`
	MajorStatus   string         `json:"majorStatus"`
	MinorStatuses []string       `json:"minorStatuses"`
	Moves         []interface{}  `json:"moves"`
	StatBoosts    map[string]int `json:"boosts"`
	Active        bool           `json:"isActive"`
	Fainted       bool           `json:"isFainted"`
	Ability       string         `json:"ability"`
	HeldItem      string         `json:"item"`
	PlayerId      string         `json:"trainerId"`
	CurrentHp     int            `json:"currentHp"`
	MaxHp         int            `json:"maxHp"`
	Shiny         bool           `json:"shiny"`
}

type BattleFieldState struct {
	Conditions []BattleFieldCondition            `json:"conditions"`
	Sides      map[string][]BattleFieldCondition `json:"sides"`
}

type BattleFieldCondition struct {
	Condition string `json:"condition"` // name of the effect
	Weather   bool   `json:"isWeather"` // flag if this condition is a weather effect
	Turns     int    `json:"turns"`     // number of turns the effect has remaining
	AltTurns  int    `json:"altTurns"`  // alternative number of turns the effect has remaining, if unknown (damp rock rain, terrain extender)
}
