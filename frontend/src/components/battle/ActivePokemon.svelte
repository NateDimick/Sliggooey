<script lang="ts">
import type { main as go } from "../../wailsjs/go/models";

export let state: go.PokemonState

function speciesForUrl(species: string): string {
    let s = species.toLowerCase()
    s = s.replaceAll(/[^a-z\-]{1}/g, "")
    return s
}
</script>

<main>
    <p>{state.nickname} Lv{state.level} {state.currentHp}/{state.maxHp}</p>
    <p>{state.majorStatus}</p>
    {#if state.shiny}
        <img src="https://play.pokemonshowdown.com/sprites/ani-shiny/{speciesForUrl(state.species)}.gif" alt="{state.nickname}">
    {:else}
        <img src="https://play.pokemonshowdown.com/sprites/ani/{speciesForUrl(state.species)}.gif" alt="{state.nickname}"> 
    {/if}
</main>