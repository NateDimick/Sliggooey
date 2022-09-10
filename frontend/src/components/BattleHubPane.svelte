<script lang="ts">
import { IPCEventTypes } from "../util";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { battleRequests, currentPaneStore, PaneInfo} from "../store"
export let info: PaneInfo

$: isFront = info.name === $currentPaneStore

EventsOn(IPCEventTypes.BattleRequest, (data) => {
    battleRequests.update((brs) => {
        brs[data.RoomId] = data
        return brs
    })
})

</script>

<main>
    <div hidden={!isFront}>
        <h1>The Battle Hub</h1>
        <div class="left-col">
            <h2>Imagine Current Games going here</h2>
        </div>
        <div class="center-col">
            <h2>Imagine Matchmaking going here</h2>
        </div>
        <div class="right-col">
            <h2>Imagine Current Challenges going here</h2>
        </div>
    </div>
</main>