<script lang="ts">
import type { main as go } from "../../wailsjs/go/models";



export let participant: go.BattleRoomParticipant

</script>

<main>
    <img src="https://play.pokemonshowdown.com/sprites/trainers/{participant.avatar}.png" alt="{participant.rating}">
    <p>{participant.name}</p>
</main>