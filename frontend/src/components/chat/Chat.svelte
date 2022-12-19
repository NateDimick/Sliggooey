<!-- a Chat component represents a PM chat with 1 user on the home pane-->
<script lang="ts">
import { pmChats } from "../../store";

import { PmSource, tsPrint } from "../../util";
import { AcceptBattleChallengeFromUser, CancelBattleChallengeToUser, RejectBattleChallengeFromUser, SendBattleChallengeToUser, SendPM } from "../../wailsjs/go/backend/App";
import ChatCommandChin from "./ChatCommandChin.svelte";
import ChatChallenge from "./ChatChallenge.svelte";

export let chatWith: string

let msgToSend: string

function sendMsg() {
    SendPM(chatWith, msgToSend)
    msgToSend = ""
}

function closeThisChat() {
    tsPrint(`closing chat box with ${chatWith}`)
    pmChats.update((pms => {
        delete pms[chatWith]
        return pms
    }))
}

</script>

<main>
    <h3>{chatWith}</h3>
    <input type="button" value="&#10006" on:click={closeThisChat}>
    {#each $pmChats[chatWith] as pm}
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