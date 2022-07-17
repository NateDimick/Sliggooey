<!-- the Chat Pane manages and contains Multiple Chat boxes -->
<script lang="ts">
import { PaneInfo, pmChats, PmRecord } from "../store";
import { IPCEventTypes, PmPayload, tsPrint, UiEventTypes } from "../util";
import { EventsEmit, EventsOn } from "../../wailsjs/runtime/runtime";
import { get } from "svelte/store";
import Chat from "./Chat.svelte";

export let info: PaneInfo

let isFront: boolean = info.front
let newChatWith: string

function openNewChat() {
    tsPrint(`new PM`)
        pmChats.update((pms) => {
            pms = [...pms, {with: newChatWith, first: {}}]
            return pms
        })
}

EventsOn(UiEventTypes.PaneChange, (paneName: string) => {
    if (paneName === info.name) {
        isFront = true
    } else {
        isFront = false
    }
})

EventsOn(IPCEventTypes.PrivateMessage, (data: PmPayload) => {
    tsPrint(`received a new PM: ${JSON.stringify(data)}`)
    let openChats = get(pmChats)
    if (!openChats.find(p => p.with === data.With)) {
        tsPrint(`new PM`)
        pmChats.update((pms: PmRecord[]) => {
            pms = [...pms, {with: data.With, first: data}]
            return pms
        })
    } else {
        // if not a new chat, then just do this:
        tsPrint(`forward chat ${JSON.stringify(data)} to chat component`)
        EventsEmit(data.With, data)
    }  
})

EventsOn(UiEventTypes.DeleteChat, (withUser: string) => {
    tsPrint(`closing chat box with ${withUser}`)
    pmChats.update((pms: PmRecord[]) => {
        pms = pms.filter((p: PmRecord) => p.with !== withUser)
        return pms
    })
})
</script>

<main>
    <div hidden={!isFront}>
        <h1>This is a Chat Pane</h1>
        <input type="text" bind:value={newChatWith}>
        <input type="button" value="Start New Chat" on:click={openNewChat}>
        {#each $pmChats as pmWith}
            <Chat chatWith={pmWith.with}/>
        {/each}
    </div>    
</main>