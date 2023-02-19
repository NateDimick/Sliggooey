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

<Trainer participant={participant}/>
<main>
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

<style>
    main {
        display: flex;
        flex-wrap: wrap;
        justify-content: center;
    }
</style>