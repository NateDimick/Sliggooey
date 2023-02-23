<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";
import ActivePokemon from "./ActivePokemon.svelte";
import Bench from "./Bench.svelte";

export let participant: go.BattleRoomParticipant
export let fieldConditions: go.BattleFieldCondition[]

function getConditions(): go.BattleFieldCondition[] {
    if (fieldConditions) {
        return fieldConditions
    } else {
        return []
    }
}

</script>

<main>
    <h2>{participant.id}</h2>
    <Bench participant={participant}/>
    {#each participant?.active as p}
        <ActivePokemon state={p}/>
    {/each}
    {#each getConditions() as condition}
        <p>{condition.condition}</p>
    {/each}
</main>

<style>
    main {
        width: 49%;
        text-align: center;
    }
</style>