import { writable } from "svelte/store";
import { PaneType, PmPayload, ViewType } from "./util";
import type { backend as go } from "./wailsjs/go/models"

// store the pokedex here
export const pokedex = writable(new Object())
// usernames as keys, list of chats as values
export const pmChats = writable(new Object())
// roomIds and keys, list of chats as values
export const roomChats = writable(new Object())
// roomIds as keys, roomStates and values
export const roomStates = writable(new Object() as {[key: string]: go.RoomState})
// list of active battle room ids (even if the pane is closed)
export const battles = writable([])
// logged in user's username
export const userName = writable("")
// list of all open Panes
export const panes = writable([
    {type: PaneType.HomePane, name: "Home", removable: false},
    {type: PaneType.RoomHubPane, name: "Chat Room Hub", removable: false},
    {type: PaneType.ChatPane, name: "Private Chats", removable: false},
    {type: PaneType.BattleHubPane, name: "Battle Hub", removable: false}
])
// name of the active pane
export const currentPaneStore = writable("Home")
// name of the active view
export const currentViewStore = writable(ViewType.Login)
// list of pending challenges
export const coldChallenges = writable([])

export type PmRecord = {
    with: string,
    first: PmPayload
}

export type PaneInfo = {
    type: PaneType,
    name: string,
    removable: boolean
}