<script lang="ts">
import BenchedPokemon from "./BenchedPokemon.svelte";
import Trainer from "./Trainer.svelte";


</script>

<main>
    <Trainer/>
    {#each ["todo"] as p}
        <BenchedPokemon/>
    {/each}
</main>