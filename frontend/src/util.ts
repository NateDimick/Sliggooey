import { LogPrint } from "../wailsjs/runtime/runtime.js";

export function tsPrint(message: string) {
    LogPrint("[TS] " + message)
}

// pane types within the client view
export enum PaneType {
    HomePane,
    ChatPane,
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

// event types that exist only within the UI
export enum UiEventTypes {
    ViewChange = "uiChangeView",
    PaneChange = "uiChangePane",
    DeletePane = "uiDeletePane",
    DeleteChat = "uiDeleteChat"
}

// event types that come from the Go backend (appTypes.go ShowdownEventTopic)
export enum IPCEventTypes {
    LoginFail = "loginFail",
    LoginSuccess = "loginSuccess",
    FormatList = "formats",
    PrivateMessage = "pm",
    RoomMessage = "roomMsg",
    Popup = "popup",
    Challenge = "challenged"
}

export enum PmSource {
    Other,
    Self,
    System
}

export type PmPayload = {
    With: string,
    IAm: String,
    Message: string,
    From: PmSource
}