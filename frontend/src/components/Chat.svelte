<!-- a Chat component represents a PM chat with 1 user on the home pane-->
<script lang="ts">
import { pmChats } from "../store";

import { PmPayload, PmSource, tsPrint } from "../util";
import { onDestroy, onMount } from "svelte";
import { get } from "svelte/store";
import { SendPM } from "../../wailsjs/go/main/App";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";

export let chatWith: string

let messages: PmPayload[] = []
let msgToSend: string

function sendMsg() {
    SendPM(chatWith, msgToSend)
    msgToSend = ""
}

function closeThisChat() {

}

function challengeToBattle() {

}

function acceptChallenge() {

}

function rejectChallenge() {

}

function cancelChallenge() {

}

function blockUser() {

}

function reportUser() {
    
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
    tsPrint("chat component destroyed")
    EventsOff(chatWith)
})
</script>

<main>
    <h3>{chatWith}</h3>
    <input type="button" value="X" on:click={closeThisChat}>
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
</main>