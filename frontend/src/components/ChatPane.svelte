<!-- the Chat Pane manages and contains Multiple Chat boxes -->
<script lang="ts">
import { coldChallenges, PaneInfo, pmChats, PmRecord } from "../store";
import { ChallengePayload, IPCEventTypes, PmPayload, tsPrint, UiEventTypes } from "../util";
import { EventsEmit, EventsOn } from "../../wailsjs/runtime/runtime";
import Chat from "./chat/Chat.svelte";

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

EventsOn(UiEventTypes.DeleteChat, (withUser: string) => {
    tsPrint(`closing chat box with ${withUser}`)
    pmChats.update((pms: PmRecord[]) => {
        pms = pms.filter((p: PmRecord) => p.with !== withUser)
        return pms
    })
})

EventsOn(IPCEventTypes.PrivateMessage, (data: PmPayload) => {
    tsPrint(`received a new PM: ${JSON.stringify(data)}`)
    if (!$pmChats.find(p => p.with === data.With)) {
        tsPrint(`new PM thread started`)
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

EventsOn(IPCEventTypes.Challenge, (data: ChallengePayload) => {
    tsPrint(`received a new Challenge: ${JSON.stringify(data)}`)
    // check if the other player is in the active pm list
    if (!$pmChats.find(p => p.with === data.With)) {
        // if place in the cold challenges store
        tsPrint("storing challenge for when component is created")
        coldChallenges.update((challenges: ChallengePayload[]) => {
            challenges = [...challenges, data]
            return challenges
        })
    } else {
        // if yes, emit event to ChatChallenge component
        tsPrint("immediately forwarding challenge to active chat")
        EventsEmit(data.With + UiEventTypes.NewChallenge, data)
    }
})

EventsOn(IPCEventTypes.ChallengeEnd, (challengeWith: string) => {
    EventsEmit(`][${challengeWith}][`, "")
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