<script lang="ts">
    import type { BattleRoomParticipant } from "../../util";
import BenchedPokemon from "./BenchedPokemon.svelte";
import Trainer from "./Trainer.svelte";

export let participant: BattleRoomParticipant

let numMysteryPokemon: number

$: {
    let numKnown = participant.inactive.length
    numMysteryPokemon = participant.teamsize - numKnown - 1
}

</script>

<main>
    <Trainer participant={participant}/>
    {#each participant.inactive as p}
        <BenchedPokemon state={p}/>
    {/each}
    {#each Array(numMysteryPokemon) as i}
        <BenchedPokemon state={null}/>
    {/each}
</main>