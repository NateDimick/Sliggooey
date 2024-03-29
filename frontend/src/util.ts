import { Log } from "./wailsjs/go/backend/App";

export function tsPrint(message: string) {
    Log(message)
}

// pane types within the client view
export enum PaneType {
    HomePane,
    ChatPane,
    BattleHubPane,
    BattlePane,
    RoomHubPane,
    RoomPane,
    TeamBuilderPane
}

// types of main application views
export enum ViewType {
    Login,
    Register,
    Client,
    Settings,
    Previous
}

// event types that come from the Go backend (appTypes.go ShowdownEventTopic)
export enum IPCEventTypes {
    LoginFail = "loginFail",
    LoginSuccess = "loginSuccess",
    FormatList = "formats",
    PrivateMessage = "pm",
    RoomInit = "newRoom",
    RoomExit = "roomExit",
    RoomMessage = "roomMsg",
    RoomState = "roomState",
    Popup = "popup",
    Games = "games",
    Challenge = "challenged",
    ChallengeEnd = "challengeEnd"
}

export enum PmSource {
    Other,
    Self,
    System
}

export type PmPayload = {
    With: string,
    IAm: string,
    Message: string,
    From: PmSource
}

export type ChallengePayload = {
    With: string,
    IAm: string,
    Format: string,
    Challenger: PmSource
}

export type NewRoomPayload = {
    RoomId: string,
    RoomType: string
}

export type RoomMessagePayload = {
    RoomId: string,
    From: string, 
    Message: string
}

export type RoomHtmlPayload = {
    RoomId: string,
    Html: string,
    Name: string,
    Update: boolean
}

/*
this is a raw example of the json included in a user login message. 
{"blockChallenges":false,"blockPMs":false,"ignoreTickets":false,"hideBattlesFromTrainerCard":false,"blockInvites":false,"doNotDisturb":false,"blockFriendRequests":false,"allowFriendNotifications":false,"displayBattlesToFriends":false,"hideLogins":false,"hiddenNextBattle":false,"inviteOnlyNextBattle":false,"language":null}
*/

export type GamesPayload = {
    games: Object,
    searching: string[]
}

export type BattleRequestPayload = {
    RequestJson: string,
    RoomId: string
}

export type BattleRequest = {
    side: PlayerSideDetails,
    rqid: number,
    active: ActivePokemon[],
    forceSwitch: boolean[],
    wait: boolean
}

export type PlayerSideDetails = {
    name: string,
    id: string,
    pokemon: SidePokemon[]
}

export type SidePokemon = {
    active: boolean,
    item: string,
    baseAbility: string,
    commanding: boolean, // omg they added this exclusively for tatsugiri
    reviving: boolean, // new in gen 9 as well, for revival blessing? why?
    ability: string,
    pokeball: string,
    moves: string[],
    stats: Map<string, number>,
    details: string,
    ident: string,
    condition: string
}

export type ActivePokemon = {
    moves: MoveInfo[],
    canMegaEvo?: boolean,    // gen 6 and 7 only, would have guessed "canMega" or "canMegaEvolve" first but noOooOoOoOOOoO
    canZMove?: MoveInfo[],   // gen 7 only, and only the name and target are provided. As many entries as moves. non-compatible moves are null
    canDynamax?: boolean,    // gen 8 only.
    maxMoves?: MaxMoveInfo,
    canTerastallize?: string // gen 9 only, string is the type the pokemon can terastallize into
    // probably a tera type field here for gen 9?
    // z moves for backwards compatibility
}

export type MaxMoveInfo = {
    maxMoves: MoveInfo[]
}

export type MoveInfo = {
    move: string,
    id: string,
    pp: number,
    maxpp: number,
    disabled: boolean,
    target: string // TODO: change to targetType once all possible values are known
}

export enum TargetType {
    Normal = "normal", // pretty sure normal = adjacent foe or ally. requires selection
    Self = "self",
    Ally = "adjacentAlly", // requires selection in triple + battles
    Allies = "allySide",
    AdjacentFoe = "adjacentFoe", // requires selection in double + battles
    Foes = "foeSide",
    AllAdjacentFoes ="allAdjacentFoes",
    All = "all",
    Random = "randomNormal"
}

export type BattleChoice = {
    Action: string,
    Move?: number,  // corresponds to the move's slot, 1-4
    AltMove?: string,
    Target?: number, // default to 0 for single battle, negative to target ally
    Gimmick?: string // max, mega, zmove, tera (? unsure of the actual string for tera ?)
}