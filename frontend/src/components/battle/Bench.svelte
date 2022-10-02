<script lang="ts">
    import type { BattleRoomParticipant } from "../../util";
import BenchedPokemon from "./BenchedPokemon.svelte";
import Trainer from "./Trainer.svelte";

export let participant: BattleRoomParticipant

</script>

<main>
    <Trainer participant={participant}/>
    {#each participant.inactive as p}
        <BenchedPokemon state={p}/>
    {/each}
</main>