<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";



export let state: go.PokemonState
</script>

<main>
    {#if state !== null}
    <!-- TODO: not every pokemon should be sliggoo on the bench -->
        <img style="width:40px; height:30px; background: url(https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png) transparent no-repeat scroll -360px -1740px;" alt="{state.nickname} : {state.species}}">
        <p>{state.nickname} : {state.species}</p>
    {:else}
        <img src="http://play.pokemonshowdown.com/sprites/itemicons/poke-ball.png" alt="Unrevealed Pokemon">
    {/if}
</main>