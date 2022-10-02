<!-- for visualizing the battle simualtion -->
<script lang="ts">
    import BattleSceneCorner from "./BattleSceneCorner.svelte";
    import { roomStates } from "../../store";
    import { RoomState, tsPrint } from "../../util";

    export let roomName: string
    let state: RoomState
    let updates: number = 0

    $: {
        state = $roomStates[roomName]
        tsPrint(`updated state in battle scene ${JSON.stringify(state)}`)
        updates += 1
    }

</script>

<main>
    <h1>{state.gen} {state.title} {state.gameType} {updates}</h1>
    {#each Object.keys(state?.participants) as p}
        <h1>{p}</h1>
        <BattleSceneCorner participant={state.participants[p]}/>
    {/each}
</main>