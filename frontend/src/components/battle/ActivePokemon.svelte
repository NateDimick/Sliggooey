<script lang="ts">
import type { PokemonState } from "../../util";

export let state: PokemonState
</script>

<main>
    <p>{state.nickname} Lv{state.level} {state.hpState}</p>
    <p>{state.majorStatus}</p>
    {#if state.shiny}
        <img src="https://play.pokemonshowdown.com/sprites/ani-shiny/{state.species.toLowerCase()}.gif" alt="{state.nickname}">
    {:else}
        <img src="https://play.pokemonshowdown.com/sprites/ani/{state.species.toLowerCase()}.gif" alt="{state.nickname}"> 
    {/if}
</main>