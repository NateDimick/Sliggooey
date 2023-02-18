<script lang="ts">
import type { backend as go } from "../../wailsjs/go/models";

export let participant: go.BattleRoomParticipant

const trainerSheet: string = "https://play.pokemonshowdown.com/sprites/trainers-sheet.png"
const spritesPerRow: number = 16
const spriteBorderEdgeLen: number = 80 // is a square

let trainerImgUrl: string
let trainerStyleStr: string = ""

let avatarNumber: number = Number(participant.avatar)
if (isNaN(avatarNumber)) {
    trainerImgUrl = `https://play.pokemonshowdown.com/sprites/trainers/${participant.avatar}.png`
} else {
    let trainerX = (avatarNumber % spritesPerRow - 1) * spriteBorderEdgeLen
    let trainerY = Math.floor(avatarNumber / spritesPerRow) * spriteBorderEdgeLen
    trainerStyleStr = `object-position: -${trainerX}px -${trainerY}px`
    trainerImgUrl = trainerSheet
}
</script>

<main>
    <img src={trainerImgUrl} style={trainerStyleStr} alt="{participant.rating}">
    <p>{participant.name}</p>
</main>

<style>
    img {
        height: 80px;
        width: 80px;
        object-fit: none;
    }
</style>