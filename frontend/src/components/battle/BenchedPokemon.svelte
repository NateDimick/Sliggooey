<script lang="ts">
import { tsPrint } from "../../util";
import { pokedex, specialIconIds } from "../../store";
import type { backend as go } from "../../wailsjs/go/models";

export let state: go.PokemonState

const spritesPerRow: number = 12
const spriteHeight: number = 30
const spriteWidth: number = 40

let pokeStyleStr: string 

$: {
    pokeStyleStr = styleSprite(state)
}

function styleSprite(state: go.PokemonState): string {
    let styling = ""
    if (state !== null) {
        let species = state?.species
        if (species !== undefined) {
            //tsPrint(`getting bench sprite for ${species}`)
            let special = species.toLowerCase().replaceAll(/[^a-z]{1}/g, "")
            //tsPrint(`special name is ${special}`)
            let dexEntry: any = $pokedex[species.toLowerCase()]
            if (dexEntry !== undefined) {
                // showdown client hard-codes all formes in a single map https://github.com/smogon/pokemon-showdown-client/blob/master/src/battle-dex-data.ts
                // which lets them do the same math below (which was reverse-engineered, for the record)
                // yuck. I hate hard coding.
                // the file can be downloaded, but it's typescript https://play.pokemonshowdown.com/src/battle-dex-data.ts
                let dexNum: number = dexEntry.num
                let pokeX = (dexNum % spritesPerRow) * spriteWidth
                let pokeY = Math.floor(dexNum / spritesPerRow) * spriteHeight
                //tsPrint(`${state.species} #${dexNum} @ -${pokeX}, -${pokeY}`)
                styling = styling + `object-position: -${pokeX}px -${pokeY}px;`
            } else if ($specialIconIds[special] !== undefined) {
                let dexNum: number = $specialIconIds[special]
                let pokeX = (dexNum % spritesPerRow) * spriteWidth
                let pokeY = Math.floor(dexNum / spritesPerRow) * spriteHeight
                //tsPrint(`${state.species} #${dexNum}! @ -${pokeX}, -${pokeY}`)
                styling = styling + `object-position: -${pokeX}px -${pokeY}px;`
            } else {
                // 0,0 is missingno/question mark icon
                styling = styling + "object-position: 0px 0px;"
            }
        }
        if (state?.isFainted) {
            styling = styling + "filter: brightness(30%);"
        }
        if (state?.isActive) {
            let brightness = 1 - state.currentHp / state.maxHp
            styling = styling + `filter: sepia(${brightness}) drop-shadow(0 0 1em goldenrod);`
        }
        if (!state?.isActive && !state?.isFainted) {
            let brightness = 1 - state.currentHp / state.maxHp
            styling = styling + `filter: sepia(${brightness});`
        }
    }
    //tsPrint(`style for ${state?.nickname} is ${styling}`)
    return styling
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