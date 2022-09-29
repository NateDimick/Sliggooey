<script lang="ts">
import { GamesPayload, IPCEventTypes, tsPrint } from "../util";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { battles } from "../store"
import CurrentGames from "./battlehub/CurrentGames.svelte";
import CurrentChallenges from "./battlehub/CurrentChallenges.svelte";
import BattleSearch from "./battlehub/BattleSearch.svelte";

EventsOn(IPCEventTypes.Games, (data: string) => {
    tsPrint(`Incoming current games info: ${data}`)
    let games: GamesPayload = JSON.parse(data)
    if (games.games !== undefined) {
        battles.set(Object.entries(games.games))
        tsPrint(`${JSON.stringify($battles)}`)
    }
    
})

</script>

<main>
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
</main>