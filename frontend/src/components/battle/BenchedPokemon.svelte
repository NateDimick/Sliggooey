<script lang="ts">
import { tsPrint } from "../../util";
import { pokedex } from "../../store";
import type { backend as go } from "../../wailsjs/go/models";

export let state: go.PokemonState

const spritesPerRow: number = 12
const spriteHeight: number = 30
const spriteWidth: number = 40

// default pokemon to ? icon

let pokeStyleStr: string 

$: {
    let species: string = state?.species
    pokeStyleStr = styleSprite(species)
}

function styleSprite(species: string): string {
    if (state !== null && state !== undefined) {
        let dexEntry: any = $pokedex[species.toLowerCase()]
        if (dexEntry !== undefined) {
            let dexNum: number = dexEntry.num
            let pokeX = (dexNum % spritesPerRow) * spriteWidth
            let pokeY = Math.floor(dexNum / spritesPerRow) * spriteHeight
            tsPrint(`${state.species} #${dexNum} @ -${pokeX}, -${pokeY}`)
            return `object-position: -${pokeX}px -${pokeY}px`
        } else {
            return "object-position: 0px 0px;"
        }
    }
}

</script>

<main>
    {#if state !== null}
        <img style={pokeStyleStr} src="https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png" alt="{state.nickname} : {state.species}}">
        <p>{state.nickname} : {state.species} :</p>
    {:else}
        <img src="http://play.pokemonshowdown.com/sprites/itemicons/poke-ball.png" alt="Unrevealed Pokemon">
    {/if}
</main>

<style>
    img {
        height: 30px;
        width: 40px;
        object-fit: none;
    }
</style>