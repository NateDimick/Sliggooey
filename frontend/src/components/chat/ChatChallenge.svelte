<!-- Display active challenge info within a private chat -->
<script lang="ts">
import { ChallengePayload, PmSource, UiEventTypes } from "../../util";
import { coldChallenges, userName } from "../../store";
import { EventsOff, EventsOn } from "../../../wailsjs/runtime/runtime";
import { AcceptBattleChallengeFromUser, CancelBattleChallengeToUser, RejectBattleChallengeFromUser } from "../../../wailsjs/go/main/App";
import { onDestroy, onMount } from "svelte";

export let challengeWith: string

let activeChallenge: boolean = false
let challenger: string = "sliggoeyAppPlaceholderChallenger"
let format: string = ""

const eventTopic = challengeWith + UiEventTypes.NewChallenge

function acceptChallenge() {
    AcceptBattleChallengeFromUser(challengeWith, "null")
    activeChallenge = false
}

function rejectChallenge() {
    RejectBattleChallengeFromUser(challengeWith)
    activeChallenge = false
}

function cancelChallenge() {
    CancelBattleChallengeToUser(challengeWith)
    activeChallenge = false
}

EventsOn(eventTopic, (data: ChallengePayload) => {
    activeChallenge = true
    format = data.Format
    if (data.Challenger === PmSource.Self) {
        challenger = $userName
    } else if (data.Challenger === PmSource.Other) {
        challenger = challengeWith
    }
})

EventsOn(`][${challengeWith}][`, () => {
    activeChallenge = false
})

onMount(() => {
    // check cold challenges if this chat was opened because of a challenge
    let c: ChallengePayload = $coldChallenges.find((c: ChallengePayload) => c.With === challengeWith)
    if (c) {
        activeChallenge = true
        format = c.Format
        if (c.Challenger === PmSource.Self) {
            challenger = $userName
        } else if (c.Challenger === PmSource.Other) {
            challenger = challengeWith
        }
    }
    // then remove the cold challenge
    coldChallenges.update((challenges: ChallengePayload[]) => {
        challenges = challenges.filter((chal: ChallengePayload) => chal.With !== challengeWith)
        return challenges
    })
})

onDestroy(() => {
    // turn off event listener
    EventsOff(eventTopic)
    EventsOff(`][${challengeWith}][`)
})
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