<script lang="ts">
import type { PaneInfo } from "../store";
import { UiEventTypes } from "../util";
import { EventsOn } from "../../wailsjs/runtime/runtime";

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
    </div>
</main>