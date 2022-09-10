<!-- the Pane Manager allows for switching between Panes and manages the addition and removal of Panes -->
<script lang="ts">
import PaneTab from "./PaneTab.svelte";
import { currentPaneStore, PaneInfo, panes, roomChats, roomStates } from "../store"
import { IPCEventTypes, NewRoomPayload, PaneType, RoomHtmlPayload, RoomMessagePayload, RoomStatePayload, tsPrint } from "../util";
import HomePane from "./HomePane.svelte";
import ChatPane from "./ChatPane.svelte";
import BattlePane from "./BattlePane.svelte";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import BattleHubPane from "./BattleHubPane.svelte";

EventsOn(IPCEventTypes.RoomInit, (data: NewRoomPayload) => {
    let newPane: PaneInfo = {name: data.RoomId, front: true, removable: true, type: undefined}
    if (data.RoomType === "chat") {
        newPane.type = PaneType.RoomPane
    }else if (data.RoomType === "battle") {
        newPane.type = PaneType.BattlePane
    }else {
        tsPrint(`Unrecognized room type: ${data.RoomType}`)
        return
    }
    roomChats.update((rms) => {
        rms[data.RoomId] = []
        return rms
    })
    panes.update((panes: PaneInfo[]) => {
        panes = [...panes, newPane]
        return panes
    })
    currentPaneStore.set(data.RoomId)
})

EventsOn(IPCEventTypes.RoomMessage, (data: RoomMessagePayload | RoomHtmlPayload) => {
    tsPrint(`received a new Room message: ${JSON.stringify(data)}`)
    tsPrint(JSON.stringify($roomChats))
    if ($roomChats[data.RoomId] === undefined) {
        // new chat
        roomChats.update((rms) => {
            rms[data.RoomId] = [data]
            return rms
        })
    } else {
        // update chat
        roomChats.update((rms) => {
            rms[data.RoomId] = [...rms[data.RoomId], data]
            return rms
        })
    }
    tsPrint(JSON.stringify($roomChats))
})

EventsOn(IPCEventTypes.RoomState, (data: RoomStatePayload) => {
    roomStates.update((rss) => {
        // Todo: need to reconcile roomstate, not just replace
        rss[data.RoomId] = data
        return rss
    })
})
</script>

<main>
    <div id="tabs">
        {#each $panes as p}
            <PaneTab info={p}/>
        {/each}
    </div>
    <div id="view">
        {#each $panes as p}
            {#if p.type === PaneType.HomePane}
                <HomePane info={p}/>
            {:else if p.type === PaneType.ChatPane}
                <ChatPane info={p}/>
            {:else if p.type === PaneType.BattleHubPane}
                <BattleHubPane info={p}/>
            {:else if p.type === PaneType.BattlePane}
                <BattlePane info={p}/>
            {/if}
        {/each}
    </div>
</main>