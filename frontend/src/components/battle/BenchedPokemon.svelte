<script lang="ts">
    import type { PokemonState } from "../../util";



export let state: PokemonState
</script>

<main>
    <!-- TODO: not every pokemon should be sliggoo on the bench -->
    <img style="width:40px; height:30px; background: url(https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png) transparent no-repeat scroll -360px -1740px;" alt="{state.nickname} : {state.species}}">
</main>