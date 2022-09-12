<!-- Display active challenge info within a private chat -->
<script lang="ts">
import { ChallengePayload, PmSource } from "../../util";
import { coldChallenges, userName } from "../../store";
import { AcceptBattleChallengeFromUser, CancelBattleChallengeToUser, RejectBattleChallengeFromUser } from "../../wailsjs/go/main/App";

export let challengeWith: string

let challenger: string = "sliggoeyAppPlaceholderChallenger"
let format: string = ""

$: activeChallenge = $coldChallenges.find((c) => c.With === challengeWith) !== undefined

coldChallenges.subscribe((challenges: ChallengePayload[]) => {
    let c = challenges.find((ch) => ch.With === challengeWith)
    if (c) {
        challenger = c.Challenger == PmSource.Self ? $userName : c.With
        format = c.Format
    }
})

function acceptChallenge() {
    AcceptBattleChallengeFromUser(challengeWith, "null")
}

function rejectChallenge() {
    RejectBattleChallengeFromUser(challengeWith)
}

function cancelChallenge() {
    CancelBattleChallengeToUser(challengeWith)
}
</script>

<main>
    <div hidden={!activeChallenge}>
        {#if challenger === challengeWith}
        <p><strong>{challenger} wants to battle!</strong></p>
        <input type="button" value="&#10004" on:click={acceptChallenge}>
        <input type="button" value="&#10006" on:click={rejectChallenge}>
        {:else}
        <p><strong>You challenged {challengeWith} to a {format} battle!</strong></p>
        <input type="button" value="Cancel" on:click={cancelChallenge}>
        {/if}
    </div>
</main>