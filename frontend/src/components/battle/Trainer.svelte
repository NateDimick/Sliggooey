<script lang="ts">
    import type { BattleRoomParticipant } from "src/util";



    export let participant: BattleRoomParticipant

</script>

<main>
    <img src="https://play.pokemonshowdown.com/sprites/trainers/{participant.avatar}.png" alt="{participant.rating}">
    <p>{participant.name}</p>
</main>