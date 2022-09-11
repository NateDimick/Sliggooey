<script lang="ts">
import { GamesPayload, IPCEventTypes, tsPrint } from "../util";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { battleRequests, battles, currentPaneStore, PaneInfo} from "../store"
import CurrentGames from "./battlehub/CurrentGames.svelte";
import CurrentChallenges from "./battlehub/CurrentChallenges.svelte";
import BattleSearch from "./battlehub/BattleSearch.svelte";
export let info: PaneInfo

$: isFront = info.name === $currentPaneStore

EventsOn(IPCEventTypes.BattleRequest, (data) => {
    tsPrint(`incoming battle request info: ${JSON.stringify(data)}`)
    battleRequests.update((brs) => {
        brs[data.RoomId] = data
        return brs
    })
    tsPrint(JSON.stringify($battleRequests))
})

EventsOn(IPCEventTypes.Games, (data: GamesPayload) => {
    tsPrint(`Incoming current games info: ${JSON.stringify(data)}`)
    if (data.games !== undefined) {
        battles.set(Object.entries(data.games))
        tsPrint(`${JSON.stringify($battles)}`)
    }
    
})

</script>

<main>
    <div hidden={!isFront}>
        <h1>The Battle Hub</h1>
        <div class="left-col">
            <CurrentGames/>
        </div>
        <div class="center-col">
            <BattleSearch/>
        </div>
        <div class="right-col">
            <CurrentChallenges/>
        </div>
    </div>
</main>