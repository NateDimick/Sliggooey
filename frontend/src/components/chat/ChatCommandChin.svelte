<!-- Provide quick buttons for Chat actions other than Chatting -->
<script lang="ts">

export let chatWith: string

import { SendBattleChallengeToUser } from "../../../wailsjs/go/main/App"

function challengeToBattle() {
    SendBattleChallengeToUser(chatWith, "gen8randombattle", "null")
}

function blockUser() {
    alert("blocking isn't supported yet - sorry")
}

function reportUser() {
    alert("Reporting isn't implemented yet - sorry")
}
</script>

<main>
    <input type="button" value="Challenge" on:click={challengeToBattle}>
    <input type="button" value="Block" on:click={blockUser}>
    <input type="button" value="Report" on:click={reportUser}>
</main>