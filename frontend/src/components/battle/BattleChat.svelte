<!-- Chat Room for a Battle Room -->
<script lang="ts">
import { roomChats } from "../../store";
import { SendRoomChat } from "../../../wailsjs/go/main/App"

export let roomName: string

let chatToSend: string = ""

function sendChat() {
    SendRoomChat(roomName, chatToSend)
    chatToSend = ""
}

</script>

<main>
    {#each $roomChats[roomName] as rm}
        {#if rm.Html === undefined}
            <p>{rm.From}: {rm.Message}</p>
        {:else}
            <div id={rm.Name}>
                {@html rm.Html}
            </div>
            {/if}
    {/each}
    <input type="text" bind:value={chatToSend}>
    <input type="button" value="Send" on:click={sendChat}>
</main>