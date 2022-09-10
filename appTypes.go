package main

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
	LoginFail          ShowdownEventTopic = "loginFail"
	LoginSuccess       ShowdownEventTopic = "loginSuccess"
	FormatTopic        ShowdownEventTopic = "formats"
	PMTopic            ShowdownEventTopic = "pm"
	NewRoomTopic       ShowdownEventTopic = "newRoom"
	RoomMessageTopic   ShowdownEventTopic = "roomMsg"
	RoomStateTopic     ShowdownEventTopic = "roomState"
	PopupTopic         ShowdownEventTopic = "popup"
	ChallengeTopic     ShowdownEventTopic = "challenged"
	ChallengeEndTopic  ShowdownEventTopic = "challengeEnd"
	BattleRequestTopic ShowdownEventTopic = "battleRequest"
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

type ShowdownUser struct {
	User       *User
	Avatar     int
	Settings   *UserSettings
	Games      *GamesStatus
	Challenges *ChallengeStatus
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
		if b.Target > 0 {
			return strings.TrimSpace(fmt.Sprintf("%s %d %d %s", b.Action, b.Move, b.Target, b.Gimmick))
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

type RoomStatePayload struct {
	RoomId       string
	Title        string
	Users        []string
	Players      []string
	GameType     string
	Gen          int
	Format       string
	IsRated      bool
	TimerOn      bool
	battleActive bool
}

type ChallengePayload struct {
	With       string
	IAm        string
	Format     string
	Challenger PmSource
}

/*
{"blockChallenges":false,"blockPMs":false,"ignoreTickets":false,"hideBattlesFromTrainerCard":false,"blockInvites":false,"doNotDisturb":false,"blockFriendRequests":false,"allowFriendNotifications":false,"displayBattlesToFriends":false,"hideLogins":false,"hiddenNextBattle":false,"inviteOnlyNextBattle":false,"language":null}
*/
type UserSettings struct {
	BlockChallenges            bool   `json:"blockChallenges"`
	BlockPMs                   bool   `json:"blockPMs"`
	IgnoreTickets              bool   `json:"ignoreTickets"`
	HideBattlesFromTrainerCard bool   `json:"hideBattlesFromTrainerCard"`
	BlockInvites               bool   `json:"blockInvites"`
	DoNotDisturb               bool   `json:"doNotDisturb"`
	BlockFriendRequests        bool   `json:"blockFriendRequests"`
	AllowFriendNotifications   bool   `json:"allowFriendNotifications"`
	DisplayBattlesToFriends    bool   `json:"displayBattlesToFriends"`
	HideLogins                 bool   `json:"hideLogins"`
	HiddenNextBattle           bool   `json:"hiddenNextBattle"`
	InviteOnlyNextBattle       bool   `json:"inviteOnlyNextBattle"`
	Language                   string `json:"language"`
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
