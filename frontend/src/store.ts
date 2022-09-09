import { writable } from "svelte/store";
import { PaneType, PmPayload } from "./util";

export const pmChats = writable(new Object())
export const roomChats = writable(new Object())
export const userName = writable("")
export const panes = writable([
    {type: PaneType.HomePane, name: "Home", front: true, removable: false}, 
    {type: PaneType.ChatPane, name: "Private Chats", front: false, removable: false}
])
export const coldChallenges = writable([])

export type PmRecord = {
    with: string,
    first: PmPayload
}

export type PaneInfo = {
    type: PaneType,
    name: string,
    front: boolean,
    removable: boolean
}