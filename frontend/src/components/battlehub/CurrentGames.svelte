<script lang="ts">
import { battles, currentPaneStore, panes, roomChats, roomStates } from "../../store";
import { JoinRoom } from "../../../wailsjs/go/main/App"
import { tsPrint } from "../../util";

function openBattleRoom(roomId: string) {
    return () => {
        // check if room is in the panes store of active panes
        if ($panes.find((p) => p.name === roomId)) {
            // if yes, just set the current Pane to roomId
            tsPrint(`switching to active battle pane ${roomId}`)
            currentPaneStore.set(roomId)
        } else {
            // if no, then need to re-join room
            tsPrint(`rejoining battle room ${roomId}`)
            JoinRoom(roomId)
        }        
    }
}
</script>

<main>
    <h1>Active Games</h1>
    {#each $battles as game}
        <button on:click={openBattleRoom(game[0])}>{game[0]} - {game[1]}</button>
    {/each}
</main>