import { LogPrint } from "../wailsjs/runtime/runtime.js";

export function tsPrint(message: string) {
    LogPrint("[TS] " + message)
}

// pane types within the client view
export enum PaneType {
    HomePane,
    ChatPane,
    BattleHubPane,
    BattlePane,
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
    ChallengeEnd = "challengeEnd",
    BattleRequest = "battleRequest"
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
    Users: string[]
}

export type GamesPayload = {
    games: Object,
    searching: string[]
}