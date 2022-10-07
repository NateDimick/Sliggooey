<script lang="ts">
import type { main as go } from "../../wailsjs/go/models";
import ActivePokemon from "./ActivePokemon.svelte";
import Bench from "./Bench.svelte";

export let participant: go.BattleRoomParticipant

</script>

<main>
    <Bench participant={participant}/>
    {#each participant?.active as p}
        <ActivePokemon state={p}/>
    {/each}
</main>