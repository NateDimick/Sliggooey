package main

import "strconv"

// types for the pokemon showdown wss protocol

type RoomType string

const (
	BattleRoom RoomType = "battle"
	ChatRoom   RoomType = "chat"
)

type MessageType string

const (
	Init             MessageType = "init"
	DeInit           MessageType = "deinit"
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
	PreviewPoke   MessageType = "poke"
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
	Drag           MessageType = "drag"
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

// p1a: Weezing
type PokemonPosition struct {
	PlayerId string
	Position int
	NickName string
}

// Weezing, L88, M
type PokemonDetails struct {
	Species string
	Level   int
	Gender  rune
	Shiny   bool
}

func NewPokemonPosition(positionSpec string) *PokemonPosition {
	p := new(PokemonPosition)
	splitSpec := NewSplitString(positionSpec, ": ")
	p.NickName = splitSpec.Get(1)
	// assume formats will only allow up to 9 players (current max is 4)
	if len(splitSpec.Get(0)) == 3 {
		p.PlayerId = splitSpec.Get(0)[:2]          // p1, p2, p3, etc but breaks for p10 or higher
		p.Position = int(splitSpec.Get(0)[2]) - 96 // 'a' is 97 in ascii/utf, we want a = 1
	} else {
		p.PlayerId = splitSpec.Get(0)
	}
	return p
}

func NewPokemonDetails(detailSpec string) *PokemonDetails {
	d := new(PokemonDetails)
	splitSpec := NewSplitString(detailSpec, ", ")
	d.Species = splitSpec.Get(0)
	d.Level = 100
	d.Shiny = false
	d.Gender = ' '
	for i := 1; i < splitSpec.len; i++ {
		token := splitSpec.Get(i)
		if token == "M" || token == "F" {
			d.Gender = rune(token[0])
		} else if token == "shiny" {
			d.Shiny = true
		} else if token[0] == 'L' {
			level, _ := strconv.Atoi(token[1:])
			d.Level = level
		}
	}
	return d
}
