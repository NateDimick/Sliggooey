<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";
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

<style>
    main {
        width: 49%;
        text-align: center;
    }
</style>