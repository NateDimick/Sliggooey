package main

import "encoding/json"

// types for the pokemon showdown wss protocol

type RoomType string

const (
	BattleRoom RoomType = "battle"
	ChatRoom   RoomType = "chat"
)

type MessageType string

const (
	Init             MessageType = "init"
	Title            MessageType = "title"
	Users            MessageType = "users"
	Message          MessageType = ""
	Html             MessageType = "html"
	Uhtml            MessageType = "uhtml"
	UhtmlChange      MessageType = "htmlchange"
	Join             MessageType = "join"
	Join2            MessageType = "j"
	Leave            MessageType = "leave"
	Leave2           MessageType = "l"
	Name             MessageType = "name"
	Name2            MessageType = "n"
	ChatMsg          MessageType = "chat"
	ChatMsg2         MessageType = "c"
	ChatTs           MessageType = "c:" // chat with a timestamp attached
	Notify           MessageType = "notify"
	BattleStart      MessageType = "battle"
	BattleStart2     MessageType = "b"
	Popup            MessageType = "popup"
	PrivateMessage   MessageType = "pm"
	UserCount        MessageType = "usercount"
	NameTaken        MessageType = "nametaken"
	ChallStr         MessageType = "challstr"
	UpdateUser       MessageType = "updateuser"
	Formats          MessageType = "formats"
	UpdateSearch     MessageType = "updatesearch"
	UpdateChallenges MessageType = "updatechallenges"
	QueryResponse    MessageType = "queryResponse"
	// battle sim message types
	Player        MessageType = "player"
	TeamSize      MessageType = "teamsize"
	GameType      MessageType = "gametype"
	Generation    MessageType = "gen"
	Tier          MessageType = "tier"
	IsRatedBattle MessageType = "rated"
	Rule          MessageType = "rule"
	ClearPoke     MessageType = "clearpoke"
	TeamPreview   MessageType = "teampreview"
	SimStart      MessageType = "start"
	Request       MessageType = "request"
	TimerOn       MessageType = "inactive"
	TimerOff      MessageType = "inactiveoff"
	Upkeep        MessageType = "upkeep"
	Turn          MessageType = "turn"
	Win           MessageType = "win"
	Tie           MessageType = "tie"
	Timestamp     MessageType = "t:"
	Timestamp2    MessageType = ":"
	// Battle sim actions (major)
	Move           MessageType = "move"
	Switch         MessageType = "switch"
	DetailsChanged MessageType = "detailschange"
	FormeChange    MessageType = "-formechange"
	Replace        MessageType = "replace"
	Swap           MessageType = "swap"
	Cannot         MessageType = "cant"
	Faint          MessageType = "faint"
	// Battle sim actions (minor)
	// move failure, no animation
	Fail     MessageType = "-fail"
	Block    MessageType = "-block"
	NoTarget MessageType = "-notarget"
	Miss     MessageType = "-miss"
	// edit pokemon hp
	Damage MessageType = "-damage"
	Heal   MessageType = "-heal"
	SetHp  MessageType = "-sethp"
	// edit pokemon status
	StatusInflict MessageType = "-status"
	StatusCure    MessageType = "-curestatus"
	TeamCure      MessageType = "-cureteam"
	EffectStart   MessageType = "-start"
	EffectEnd     MessageType = "-end"
	// edit pokemon stat boosts
	Boost         MessageType = "-boost"
	Unboost       MessageType = "-unboost"
	SetBoost      MessageType = "-setboost"
	SwapBoost     MessageType = "-swapboost"
	InvertBoost   MessageType = "-invertboost"
	ClearBoost    MessageType = "-clearboost"
	ClearAllBoost MessageType = "-clearallboost"
	ClearPosBoost MessageType = "-clearpositiveboost"
	ClearNegBoost MessageType = "-clearnegativeboost"
	CopyBoost     MessageType = "-copyBoost"
	// field conditions
	Weather    MessageType = "-weather"
	FieldStart MessageType = "-fieldstart"
	FieldEnd   MessageType = "-fieldend"
	SideStart  MessageType = "-sidestart"
	SideEnd    MessageType = "-sideend"
	// damage efficacy
	Critical  MessageType = "-crit"
	Effective MessageType = "-supereffective"
	Resisted  MessageType = "-resisted"
	Immune    MessageType = "-immune"
	// misc events
	Item             MessageType = "-item"
	ItemEnd          MessageType = "-enditem"
	Ability          MessageType = "-ability"
	AbilityEnd       MessageType = "-endability"
	MegaEvolve       MessageType = "-mega"
	PrimalForm       MessageType = "-primal"
	Burst            MessageType = "-burst"
	ZMove            MessageType = "-zpower"
	ZProtect         MessageType = "-zbroken"
	ActivateEffect   MessageType = "-activate"
	Hint             MessageType = "-hint"
	Center           MessageType = "-center"
	MiscMsg          MessageType = "-message"
	CombineMove      MessageType = "-combine"
	WaitForMove      MessageType = "-waiting"
	PrepareMove      MessageType = "-prepare"
	Recharge         MessageType = "-mustrecharge"
	Nothing          MessageType = "-nothing" // deprecated, but if it's listed then I'll support it
	HitCount         MessageType = "-hitcount"
	SingleMoveEffect MessageType = "-singlemove"
	SingleTurnEffect MessageType = "-singleturn"
)

type ChatCommand string

const (
	// all commands listed by /help and grouped the same way
	// regular commands
	Report    ChatCommand = "/report"
	Reply     ChatCommand = "/reply"
	Logout    ChatCommand = "/logout"
	Challenge ChatCommand = "/challenge"
	Search    ChatCommand = "/search"
	Rating    ChatCommand = "/rating"
	WhoIs     ChatCommand = "/whois"
	UserCmd   ChatCommand = "/user" // help says does not exist
	JoinCmd   ChatCommand = "/join"
	LeaveCmd  ChatCommand = "/leave"
	UserAuth  ChatCommand = "/userauth"
	RoomAuth  ChatCommand = "/roomauth"

	// battle room commands
	SaveReplay  ChatCommand = "/savereplay"
	HideRoom    ChatCommand = "/hideroom"
	InviteOnly  ChatCommand = "/inviteonly"
	InviteUser  ChatCommand = "/invite"
	TimerToggle ChatCommand = "/timer"
	Forfeit     ChatCommand = "/forfeit"

	// option commands
	Nick            ChatCommand = "/nick" // no help info. may be deprecated
	AvatarChange    ChatCommand = "/avatar"
	IgnoreUser      ChatCommand = "/ignore"
	SetStatus       ChatCommand = "/status"
	ClearStatus     ChatCommand = "/clearstatus"
	AmAway          ChatCommand = "/away"
	AmBusy          ChatCommand = "/busy"
	AmBack          ChatCommand = "/back"
	DoNotDisturb    ChatCommand = "/donotdisturb" // like /busy but also blocks pms and challenges
	Timestamps      ChatCommand = "/timestamps"   // reportedly not a command anymore, likely deprecated but not removed from /help's list
	Highlight       ChatCommand = "/highlight"
	ShowJoins       ChatCommand = "/showjoins"
	HideJoins       ChatCommand = "/hidejoins"
	BlockChallenges ChatCommand = "/blockchallenges"
	BlockPms        ChatCommand = "/blockpms"

	// information commands
	// can be "broadcast" by replacing / with !, but requires rank
	Groups      ChatCommand = "/groups"
	Faq         ChatCommand = "/faq"
	Rules       ChatCommand = "/rules"
	Intro       ChatCommand = "/intro"
	FormatsHelp ChatCommand = "/formatshelp"
	OtherMetas  ChatCommand = "/othermetas" // can also be /om
	Analysis    ChatCommand = "/analysis"
	Punishments ChatCommand = "/punishments"
	Calculator  ChatCommand = "/calc" // /rcalc, /bsscalc as well. provides links to calculators
	Git         ChatCommand = "/git"  // /opensource as well. links to PSS repos
	Cap         ChatCommand = "/cap"  // "create a pokemon" info
	RoomHelp    ChatCommand = "/roomhelp"
	Help        ChatCommand = "/help" // /h /? all the same
	RoomFaq     ChatCommand = "/roomfaq"

	// data commands
	// many of these can be broadcast
	Data          ChatCommand = "/data"
	DexSearch     ChatCommand = "/dexsearch"  // /ds
	MoveSearch    ChatCommand = "/movesearch" // /ms
	ItemSearch    ChatCommand = "/itemsearch" // /is
	Learn         ChatCommand = "/learn"
	StatCalc      ChatCommand = "/statcalc"
	Effectiveness ChatCommand = "/effectiveness"
	Weakness      ChatCommand = "/weakness"
	Coverage      ChatCommand = "/coverage"
	RandomMove    ChatCommand = "/randommove"
	RandomPokemon ChatCommand = "/randompokemon"

	// not in the /help list
	UseTeam         ChatCommand = "/utm"
	AcceptChallenge ChatCommand = "/accept"
	RejectChallenge ChatCommand = "/reject"
	CancelChallenge ChatCommand = "/cancelchallenge"
	CancelSearch    ChatCommand = "/cancelsearch"
	Choose          ChatCommand = "/choose"
	Rename          ChatCommand = "/trn"
	Log             ChatCommand = "/log"      // not actually a command, but will come in sometimes from server
	Me              ChatCommand = "/me"       // attributes action to user
	Announce        ChatCommand = "/announce" // or /wall requires rank
	CmdError        ChatCommand = "/error"
	Text            ChatCommand = "/text"
	Raw             ChatCommand = "/raw"
	NoNotify        ChatCommand = "/nonotify"
	SendMessage     ChatCommand = "/msg"
)

type QueryType string

const (
	RoomList    QueryType = "roomlist"
	UserDetails QueryType = "userdetails"
)

type UserRank rune

const (
	//https://www.smogon.com/forums/threads/pok%C3%A9mon-showdown-forum-rules-resources-read-here-first.3570628/#post-6774482
	NoRank    UserRank = ' '
	Voice     UserRank = '+'
	Driver    UserRank = '%'
	Moderator UserRank = '@'
	Owner     UserRank = '#'
	Admin     UserRank = '&'
)

type LoginResponse struct {
	LoggedIn      bool   `json:"loggedin"`
	UserName      string `json:"username"`
	Assertion     string `json:"assertion"`
	ActionSuccess bool   `json:"actionsuccess"`
}

type GamesStatus struct {
	// from the updatesearch status message
	Games  json.RawMessage `json:"games"`
	Search []string        `json:"searching"`
}

func (g *GamesStatus) GetGames() map[string]string {
	games := make(map[string]string)
	if string(g.Games) != "null" {
		json.Unmarshal(g.Games, &games)
	}
	return games
}

type ChallengeStatus struct {
	// from the updatechallenges status message
	Challengers map[string]string `json:"challengesFrom"`
	Challenge   *ChallengeTo      `json:"challengeTo"`
}

type ChallengeTo struct {
	Opponent string `json:"to"`
	Format   string `json:"format"`
}

type FormatInfo []FormatSection

type FormatSection map[string][]Format

type Format struct {
	Name string
	Tags string
}

type BattleRequest struct {
	Side      *PlayerSideDetails `json:"side"`
	Active    []ActiveDetails    `json:"active"`
	RequestId int                `json:"rqid"`
}

type PlayerSideDetails struct {
	PlayerName     string        `json:"name"`
	PlayerId       int           `json:"id"`
	PlayersPokemon []SidePokemon `json:"pokemon"`
}

type SidePokemon struct {
	Active      bool           `json:"active"`
	Item        string         `json:"item"`
	BaseAbility string         `json:"baseAbility"`
	Ability     string         `json:"ability"`
	Pokeball    string         `json:"pokeball"`
	Moves       []string       `json:"moves"`
	Stats       map[string]int `json:"stats"`
	Details     string         `json:"details"`
	Identifier  string         `json:"ident"`
	Condition   string         `json:"condition"`
}

type ActiveDetails struct {
	Moves      []MoveInfo `json:"moves"`
	CanDynamax bool       `json:"canDynamax"`
	DmaxMoves  []MoveInfo `json:"maxMoves"`
}

type MoveInfo struct {
	MoveName  string `json:"move"`
	MoveId    string `json:"id"`
	CurrentPP int    `json:"pp"`
	MaxPP     int    `json:"maxpp"`
	Disabled  bool   `json:"disabled"`
	Target    string `json:"target"`
}
