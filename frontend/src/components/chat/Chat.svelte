<!-- a Chat component represents a PM chat with 1 user on the home pane-->
<script lang="ts">
import { pmChats } from "../../store";

import { PmPayload, PmSource, tsPrint, UiEventTypes } from "../../util";
import { onDestroy, onMount } from "svelte";
import { get } from "svelte/store";
import { AcceptBattleChallengeFromUser, CancelBattleChallengeToUser, RejectBattleChallengeFromUser, SendBattleChallengeToUser, SendPM } from "../../../wailsjs/go/main/App";
import { EventsEmit, EventsOff, EventsOn } from "../../../wailsjs/runtime/runtime";
import ChatCommandChin from "./ChatCommandChin.svelte";
import ChatChallenge from "./ChatChallenge.svelte";

export let chatWith: string

let messages: PmPayload[] = []
let msgToSend: string

function sendMsg() {
    SendPM(chatWith, msgToSend)
    msgToSend = ""
}

function closeThisChat() {
    EventsEmit(UiEventTypes.DeleteChat, chatWith)
}

EventsOn(chatWith, (data: PmPayload) => {
    tsPrint(`chat component with ${chatWith} received new message ${JSON.stringify(data)}`)
    messages = [...messages, data]
    tsPrint(`chat history: ${JSON.stringify(messages)}`)
})

onMount(() => {
    // get first message from store
    tsPrint(`New Chat component with user ${chatWith}`)
    messages = [get(pmChats).find(p => p.with === chatWith).first]
})

onDestroy(() => {
    tsPrint(`chat component with user ${chatWith} destroyed`)
    EventsOff(chatWith)
})
</script>

<main>
    <h3>{chatWith}</h3>
    <input type="button" value="&#10006" on:click={closeThisChat}>
    {#each messages as pm}
        {#if pm.From === PmSource.Other}
            <p>{pm.With} {pm.Message}</p>
        {:else if pm.From === PmSource.Self}
            <p>{pm.IAm} {pm.Message}</p>
        {:else if pm.From === PmSource.System}
            <p>{pm.Message}</p>
        {/if}
    {/each}
    <input type="text" bind:value={msgToSend}>
    <input type="button" value="Send" on:click={sendMsg}>
    <ChatChallenge challengeWith={chatWith}/>
    <ChatCommandChin {chatWith}/>
</main>