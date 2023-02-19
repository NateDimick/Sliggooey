<script lang="ts">
import { tsPrint } from "../../util";
import { pokedex, specialIconIds } from "../../store";
import type { backend as go } from "../../wailsjs/go/models";

export let state: go.PokemonState

const spritesPerRow: number = 12
const spriteHeight: number = 30
const spriteWidth: number = 40

// default pokemon to ? icon

let pokeStyleStr: string 

$: {
    let species: string = state?.species
    if (species !== undefined) {
        pokeStyleStr = styleSprite(species)
    }
}

function styleSprite(species: string): string {
    tsPrint(`getting bench sprite for ${species}`)
    let special = species.toLowerCase().replaceAll(/[^a-z]{1}/g, "")
    tsPrint(`special name is ${special}`)
    let dexEntry: any = $pokedex[species.toLowerCase()]
    if (dexEntry !== undefined) {
        // showdown client hard-codes all formes in a single map https://github.com/smogon/pokemon-showdown-client/blob/master/src/battle-dex-data.ts
        // which lets them do the same math below (which was reverse-engineered, for the record)
        // yuck. I hate hard coding.
        // the file can be downloaded, but it's typescript https://play.pokemonshowdown.com/src/battle-dex-data.ts
        let dexNum: number = dexEntry.num
        let pokeX = (dexNum % spritesPerRow) * spriteWidth
        let pokeY = Math.floor(dexNum / spritesPerRow) * spriteHeight
        tsPrint(`${state.species} #${dexNum} @ -${pokeX}, -${pokeY}`)
        return `object-position: -${pokeX}px -${pokeY}px`
    } else if ($specialIconIds[special] !== undefined) {
        let dexNum: number = $specialIconIds[special]
        let pokeX = (dexNum % spritesPerRow) * spriteWidth
        let pokeY = Math.floor(dexNum / spritesPerRow) * spriteHeight
        tsPrint(`${state.species} #${dexNum}! @ -${pokeX}, -${pokeY}`)
        return `object-position: -${pokeX}px -${pokeY}px`
    } else {
        return "object-position: 0px 0px;"
    }
}

</script>

<main>
    {#if state !== null}
        <!-- [Gen9] if bounsweet/steenee appear, that means that object-position style was not set-->
        <img style={pokeStyleStr} src="https://play.pokemonshowdown.com/sprites/pokemonicons-sheet.png" alt="{state.nickname} : {state.species}}">
        <!-- <p>{state.nickname} : {state.species} :</p> -->
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