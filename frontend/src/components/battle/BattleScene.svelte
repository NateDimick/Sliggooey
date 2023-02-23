<!-- for visualizing the battle simualtion -->
<script lang="ts">
import BattleSceneCorner from "./BattleSceneCorner.svelte";
import { roomStates } from "../../store";
import type { backend as go } from "../../wailsjs/go/models";

export let roomName: string
let state: go.RoomState
let updates: number = 0

$: {
    state = $roomStates[roomName]
    // tsPrint(`updated state in battle scene ${JSON.stringify(state)}`)  // if command pallette ever eats shit, uncomment this line
    updates += 1
}

function getConditions() {
    if (state?.field?.conditions) {
        return state.field.conditions
    } else {
        return []
    }
}

</script>

<main>
    <h1>gen {state.gen} {state.title} {state.gameType} {updates}</h1>
    {#each getConditions() as condition}
        <!-- Perhaps this should be split out in to it's own component later when dealing with graphics/etc -->
        <p>{condition.condition} {condition.turns} or {condition.altTurns}</p>
    {/each}
    <div class="battle-field">
        {#each Object.keys(state?.participants) as p}
            <BattleSceneCorner participant={state.participants[p]} fieldConditions={state?.field?.sides[p]}/>
        {/each}
    </div>
</main>

<style>
    .battle-field {
        display: flex;
    }
</style>