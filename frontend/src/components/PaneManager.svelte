<!-- the Pane Manager allows for switching between Panes and manages the addition and removal of Panes -->
<script lang="ts">
import PaneTab from "./PaneTab.svelte";
import { currentPaneStore, PaneInfo, panes, roomChats, roomStates } from "../store"
import { IPCEventTypes, NewRoomPayload, PaneType, RoomHtmlPayload, RoomMessagePayload, tsPrint } from "../util";
import HomePane from "./HomePane.svelte";
import ChatPane from "./ChatPane.svelte";
import BattlePane from "./BattlePane.svelte";
import { EventsOn } from "../wailsjs/runtime/runtime";
import BattleHubPane from "./BattleHubPane.svelte";
import RoomsHubPane from "./RoomsHubPane.svelte";
import RoomPane from "./RoomPane.svelte";
import { backend as go } from "../wailsjs/go/models"
import { ReconcileRoomState } from "../wailsjs/go/backend/App";
import { Mutex } from 'async-mutex'

const stateUpdateMutex = new Mutex()

EventsOn(IPCEventTypes.RoomInit, (data: NewRoomPayload) => {
    let newPane: PaneInfo = {name: data.RoomId, removable: true, type: undefined}
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
    roomStates.update((rss) => {
        let newState = new go.RoomState({title: "brand new battle room", participants: {}, request: null})
        tsPrint(`new room state: ${printRoomState(newState)}`)
        rss[data.RoomId] = newState
        return rss
    })
    panes.update((panes: PaneInfo[]) => {
        panes = [...panes, newPane]
        return panes
    })
    currentPaneStore.set(data.RoomId)
})

EventsOn(IPCEventTypes.RoomExit, (data: string) => {
    tsPrint(`closing room ${data}`)
    roomChats.update((rms) => {
        delete rms[data]
        return rms
    })
    roomStates.update((rss) => {
        delete rss[data]
        return rss
    })
    panes.update((panes: PaneInfo[]) => {
        panes = panes.filter((p) => p.name !== data)
        return panes
    })
    if ($currentPaneStore === data) {
        currentPaneStore.set("Home")
    }
})

EventsOn(IPCEventTypes.RoomMessage, (data: RoomMessagePayload | RoomHtmlPayload) => {
    tsPrint(`received a new Room message: ${JSON.stringify(data)}`)
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
})

EventsOn(IPCEventTypes.RoomState, async (data) => {
    //let updateData = go.UpdateRoomStatePayload.createFrom(data)
    if($roomStates[data.RoomId]) {
        tsPrint(`Updating room state: ${data.RoomId}`)
        // must wrap state updating in a mutex to ensure only one update occurs at a time
        await stateUpdateMutex.runExclusive(async () => {
            let updatedState = await ReconcileRoomState(data, $roomStates[data.RoomId])
            tsPrint(`State after update: ${data.RoomId}, ${printRoomState(updatedState)}`)
            roomStates.update((rss) => {
                rss[data.RoomId] = updatedState
                rss = rss
                return rss
            })
        }) 
    } else {
        tsPrint(`Room state update ${data.RoomId} received but user is not in that room`)
    }
    
})

function printRoomState(s: go.RoomState): string {
    let x = {
        request: s.request,
        title: s.title,
        gen: s.gen,
        gametype: s.gameType,
        tier: s.tier,
        players: s.participants
    }
    return JSON.stringify(x)
}

function printRoomStateUpdate(u: go.UpdateRoomStatePayload): string {
    let x = {
        roomid: u.RoomId,
        rated: u.rated,
        request: u.request,
        gen: u.gen,
        gametype: u.gameType, 
        tier: u.tier,
        playerId: u.player.playerId,
        avatar: u.player.avatarName,
        pspecies: u.player.activePokemonUpdate.positionalDetails
    }

    return JSON.stringify(x)
}
</script>

<main>
    <div id="tabs">
        {#each $panes as p}
            <PaneTab info={p}/>
        {/each}
    </div>
    <div id="view">
        {#each $panes as p}
        <div hidden={p.name !== $currentPaneStore}>
            {#if p.type === PaneType.HomePane}
                <HomePane/>
            {:else if p.type === PaneType.ChatPane}
                <ChatPane/>
            {:else if p.type === PaneType.BattleHubPane}
                <BattleHubPane/>
            {:else if p.type === PaneType.BattlePane}
                <BattlePane info={p}/>
            {:else if p.type === PaneType.RoomHubPane}
                <RoomsHubPane/>
            {:else if p.type === PaneType.RoomPane}
                <RoomPane info={p}/>
            {/if}
        </div>
        {/each}
    </div>
</main>