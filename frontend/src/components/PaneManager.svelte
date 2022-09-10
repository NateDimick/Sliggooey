<!-- the Pane Manager allows for switching between Panes and manages the addition and removal of Panes -->
<script lang="ts">
import PaneTab from "./PaneTab.svelte";
import { PaneInfo, panes, roomChats } from "../store"
import { IPCEventTypes, NewRoomPayload, PaneType, RoomHtmlPayload, RoomMessagePayload, tsPrint } from "../util";
import HomePane from "./HomePane.svelte";
import ChatPane from "./ChatPane.svelte";
import BattlePane from "./BattlePane.svelte";
import { EventsOn } from "../../wailsjs/runtime/runtime";

EventsOn(IPCEventTypes.RoomInit, (data: NewRoomPayload) => {
    let newPane: PaneInfo = {name: data.RoomId, front: false, removable: true, type: undefined}
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
            {:else if p.type === PaneType.BattlePane}
                <BattlePane info={p}/>
            {/if}
        {/each}
    </div>
</main>