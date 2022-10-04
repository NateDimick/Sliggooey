import { LogPrint } from "./wailsjs/runtime/runtime";

export function tsPrint(message: string) {
    LogPrint("[TS] " + message)
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

export type RoomStatePayload = {
    RoomId: string,
    Title: string,
    // only for battle rooms
    Gen: number,
    Tier: string,
    GameType: string,
    Timer: boolean,
    Rated: boolean,
    Active: boolean,
    Player: PlayerPayload,
    Request: string
}

export type RoomState = {
    request: BattleRequest,
    title: string,
    gen: number,
    tier: string,
    gameType: string, // TODO: may become enum
    timer: boolean,
    rated: boolean,
    active: boolean,
    participants: Object // keys are player ids, values are BattleRoomParticipant
}

export type BattleRoomParticipant = {
    name: string,
    id: string,
    rating: string,
    teamsize: number,
    avatar: string,
    active: PokemonState[],
    inactive: PokemonState[]
}

export type GamesPayload = {
    games: Object,
    searching: string[]
}

export type BattleRequestPayload = {
    RequestJson: string,
    RoomId: string
}

export type PlayerPayload = {
    PlayerId: string,
    Name: string, 
    Avatar: string,
    Rating: string,
    TeamSize: number,
    ActivePokemon: PlayerActivePokemonPayload
}

export type PlayerActivePokemonPayload = {
    Reason: string,
    Position,
    Details,
    HP,
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
    canDynamax?: boolean,
    maxMoves?: MaxMoveInfo
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
    Normal = "normal",
    Self = "self",
    Ally = "adjacentAlly",
    Allies = "allySide",
    AdjacentFoe = "adjacentFoe",
    Foes = "foeSide",
    All = "all",
    Random = "randomNormal"
}

export type BattleChoice = {
    Action: string,
    Move?: number,  // corresponds to the move's slot, 1-4
    AltMove?: string,
    Target?: number, // default to 0 for single battle, negative to target ally
    Gimmick?: string // max, mega, zmove
}

export type PokemonState = {
    species?: string,
    nickname?: string,
    gender?: string,
    level?: number,
    majorStatus?: string,
    minorStatuses?: string[],
    moves?: MoveInfo[],
    statBoosts?: number[],
    active?: boolean,
    fainted?: boolean,
    playerId?: string,
    currentHp?: number,
    maxHp?: number,
    shiny?: boolean
}

export function newRoomState(): RoomState {
    return {
        request: null,
        title: "",
        gen: null,
        tier: "",
        gameType: "",
        timer: false,
        rated: false,
        active: false,
        participants: new Object()
    }
}

export function reconcileRoomState(update: RoomStatePayload, base: RoomState): RoomState {
    if (update.Title) {
        base.title = update.Title
    }
    if (update.Gen) {
        base.gen = update.Gen
    }
    if (update.GameType) {
        base.gameType = update.GameType
    }
    if (update.Request !== "" && update.Request !== null) {
        //tsPrint("Current request updated")
        base.request = JSON.parse(update.Request)
    }
    if (update.Tier) {
        base.tier = update.Tier
    }
    if (update.Rated) {
        base.rated = update.Rated
    }
    if (update.Timer !== undefined) {
        base.timer = update.Timer
    }
    if (update.Active !== undefined) {
        base.active = update.Active
    }
    if (update.Player) {
        if (base.participants[update.Player.PlayerId]) {
            base.participants[update.Player.PlayerId] = reconcilePlayerState(update.Player, base.participants[update.Player.PlayerId])
        } else {
            base.participants[update.Player.PlayerId] = {
                name: update.Player.Name,
                avatar: update.Player.Avatar,
                elo: update.Player.Rating,
                teamsize: 0,
                id: update.Player.PlayerId,
                active: [],
                inactive: []
            }
        }
    } else {
        tsPrint(`No Player or player pokemon update ${update}`)
    }

    return base
}

function reconcilePlayerState(update: PlayerPayload, base: BattleRoomParticipant): BattleRoomParticipant {
    if (update.TeamSize) {
        base.teamsize = update.TeamSize
    }
    if (update.ActivePokemon) {
        tsPrint(`Updating active pokemon ${JSON.stringify(update)}`)
        let newActivePokemon = update.ActivePokemon
        let newState: PokemonState = {
            species: newActivePokemon.Details?.Species,
            level: newActivePokemon.Details?.Level,
            gender: newActivePokemon.Details?.Gender,
            shiny: newActivePokemon.Details?.Shiny,
            nickname: newActivePokemon.Position?.NickName,
            playerId: newActivePokemon.Position?.PlayerId,
            currentHp: newActivePokemon.HP?.Current,
            maxHp: newActivePokemon.HP?.Max,
            majorStatus: newActivePokemon.HP?.Status,
            minorStatuses: [],
            moves: [],
            statBoosts: [],
            active: true,
            fainted: false
        }
        while (newActivePokemon.Position.Position > base.active.length) {
            base.active.push(null)
        }
        if (base.active[newActivePokemon.Position.Position - 1] && (newActivePokemon.Reason === "switch" || newActivePokemon.Reason === "drag")) {
            let switchingOut = base.active[newActivePokemon.Position.Position - 1]
            switchingOut.active = false
            base.inactive = [...base.inactive, switchingOut]  // put formerly active pokemon back in inactive list
            base.inactive = base.inactive.filter((p: PokemonState) => {return !equalPokeState(p, newState)}) // remove new active pokemon from inactive list
            base.active[newActivePokemon.Position.Position - 1] = newState // put newly active pokemon on the active list
        } else if (newActivePokemon.Reason === "switch" && base.active[newActivePokemon.Position.Position - 1] === null) {
            base.active[newActivePokemon.Position.Position - 1] = newState // put newly active pokemon on the active list
        }
        if (newActivePokemon.Reason === "faint") {
            newState.fainted = true
            tsPrint(`${newState.nickname} fainted`)
        } else if (newActivePokemon.Reason === "hpupdate") {
            let active = base.active[newActivePokemon.Position.Position - 1]
            active.currentHp = newState.currentHp
            active.maxHp = newState.maxHp
            active.majorStatus = newState.majorStatus
            base.active[newActivePokemon.Position.Position - 1] = active
            tsPrint(`${newState.nickname} hp changed`)
        } else {
            tsPrint(`Not a recognized update reason: ${newActivePokemon.Reason}`)
        }
        base.active = [...base.active]
        tsPrint(`updated active pokemon: [${newActivePokemon.Reason}] ${JSON.stringify(base.active)} | ${JSON.stringify(base.inactive)}`)
    } else {
        tsPrint(`No player pokemon update ${JSON.stringify(update)}`)
    }
    return base
}

function equalPokeState(a: PokemonState, b: PokemonState): boolean {
    return a.level === b.level 
        && a.species === b.species 
        && a.shiny === b.shiny 
        && a.nickname === b.nickname 
        && a.maxHp === b.maxHp 
        && a.gender === b.gender
}