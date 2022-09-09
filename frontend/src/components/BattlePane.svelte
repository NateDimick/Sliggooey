<script lang="ts">
import type { PaneInfo } from "../store";
import { UiEventTypes } from "../util";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import BattleScene from "./battle/BattleScene.svelte";
import BattleCommandPalette from "./battle/BattleCommandPalette.svelte";
import BattleChat from "./battle/BattleChat.svelte";

export let info: PaneInfo

let isFront: boolean = info.front

EventsOn(UiEventTypes.PaneChange, (paneName: string) => {
    if (paneName === info.name) {
        isFront = true
    } else {
        isFront = false
    }
})
</script>

<main>
    <div hidden={!isFront}>
        <h1>This is a battle pane</h1>
        <div id="left-col">
            <h2>Imagine something cool being here, like a damage calculator</h2>
        </div>
        <div id="center-col">
            <BattleScene/>
            <BattleCommandPalette/>
        </div>
        <div id="right-col">
            <BattleChat roomName={info.name}/>
        </div>
    </div>
</main>