<script lang="ts">
    import type { BattleRoomParticipant } from "../../util";
import ActivePokemon from "./ActivePokemon.svelte";
import Bench from "./Bench.svelte";

export let participant: BattleRoomParticipant

</script>

<main>
    <Bench participant={participant}/>
    {#each participant?.active as p}
        <ActivePokemon state={p}/>
    {/each}
</main>