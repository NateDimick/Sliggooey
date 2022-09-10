<!-- the Chat Pane manages and contains Multiple Chat boxes -->
<script lang="ts">
import { coldChallenges, PaneInfo, pmChats, currentPaneStore } from "../store";
import { ChallengePayload, IPCEventTypes, PmPayload, tsPrint } from "../util";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import Chat from "./chat/Chat.svelte";

export let info: PaneInfo

$: isFront = info.name === $currentPaneStore
let newChatWith: string

function openNewChat() {
    tsPrint(`new PM`)
        pmChats.update((pms) => {
            pms[newChatWith] = []
            return pms
        })
}

EventsOn(IPCEventTypes.PrivateMessage, (data: PmPayload) => {
    tsPrint(`received a new PM: ${JSON.stringify(data)}`)
    tsPrint(JSON.stringify($pmChats))
    if ($pmChats[data.With] === undefined) {
        // new chat
        pmChats.update((pms) => {
            pms[data.With] = [data]
            return pms
        })
    } else {
        // update chat
        pmChats.update((pms) => {
            pms[data.With] = [...pms[data.With], data]
            return pms
        })
    }
    tsPrint(JSON.stringify($pmChats))
})

EventsOn(IPCEventTypes.Challenge, (data: ChallengePayload) => {
    tsPrint(`received a new Challenge: ${JSON.stringify(data)}`)
    // check if the other player is in the active pm list
    coldChallenges.update((challenges: ChallengePayload[]) => {
        challenges = [...challenges, data]
        return challenges
    })
})

EventsOn(IPCEventTypes.ChallengeEnd, (challengeWith: string) => {
    coldChallenges.update((challenges: ChallengePayload[]) => {
        challenges = challenges.filter((c) => c.With !== challengeWith)
        return challenges
    })
})
</script>

<main>
    <div hidden={!isFront}>
        <h1>This is a Chat Pane</h1>
        <input type="text" bind:value={newChatWith}>
        <input type="button" value="Start New Chat" on:click={openNewChat}>
        {#each Object.keys($pmChats) as pmWith}
            <Chat chatWith={pmWith}/>
        {/each}
    </div>    
</main>