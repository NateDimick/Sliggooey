<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";

export let state: go.PokemonState
let imageName: string = speciesForUrl(state.species)
let sprite = new Image()
let fallbackLevel: number = 0

sprite.onerror = () => {
    fallbackLevel++
    sprite.src = getUrl(state, fallbackLevel)
}

function speciesForUrl(species: string): string {
    let s = species.toLowerCase()
    s = s.replaceAll(/[^a-z\-]{1}/g, "")
    // need to remove the second hyphen and the second hypen ONLY
    return s
}

function getUrl(pokemon: go.PokemonState, fallback: number): string {
    let species = speciesForUrl(pokemon.species)
    let options = {
        regular: [
            `https://play.pokemonshowdown.com/sprites/ani/${species}.gif`,
            `https://play.pokemonshowdown.com/sprites/gen5/${species}.png`,
            "https://play.pokemonshowdown.com/sprites/gen5/0.png"
        ],
        shiny: [
            `https://play.pokemonshowdown.com/sprites/ani-shiny/${species}.gif`,
            `https://play.pokemonshowdown.com/sprites/gen5-shiny/${species}.png`,
            "https://play.pokemonshowdown.com/sprites/gen5/0.png"
        ]
    }
    let listToCheck = pokemon.shiny ? options.shiny : options.regular
    return listToCheck[fallback]
}

sprite.src = getUrl(state, fallbackLevel)
$: {
    if (speciesForUrl(state.species) !== imageName) {
        imageName = speciesForUrl(state.species)
        fallbackLevel = 0
        sprite.src = getUrl(state, fallbackLevel)
    }
}
</script>

<main>
    <p>{state.nickname} Lv{state.level} {state.currentHp}/{state.maxHp}</p>
    <p>{state.majorStatus}</p>
    {#each state.minorStatuses as ms}
        <p>{ms}</p>
    {/each}
    {#each Object.keys(state.boosts) as boost}
        <p>{boost} {state.boosts[boost]}</p>
    {/each}
    <img src={sprite.src} alt="{state.nickname}">
</main>