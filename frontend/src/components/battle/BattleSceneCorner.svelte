<script lang="ts">
import ActivePokemon from "./ActivePokemon.svelte";
import Bench from "./Bench.svelte";

</script>

<main>
    <Bench/>
    {#each ["todo"] as p}
        <ActivePokemon/>
    {/each}
</main>