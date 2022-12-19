<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";
import BenchedPokemon from "./BenchedPokemon.svelte";
import Trainer from "./Trainer.svelte";

export let participant: go.BattleRoomParticipant

let numMysteryPokemon: number

$: {
    let numKnown = participant.inactive.length
    numMysteryPokemon = Math.max(participant.teamSize - numKnown - 1, 0)
}

</script>

<main>
    <Trainer participant={participant}/>
    {#each participant.active as p}
        <BenchedPokemon state={p}/>
    {/each}
    {#each participant.inactive as p}
        <BenchedPokemon state={p}/>
    {/each}
    {#each Array(numMysteryPokemon) as i}
        <BenchedPokemon state={null}/>
    {/each}
</main>