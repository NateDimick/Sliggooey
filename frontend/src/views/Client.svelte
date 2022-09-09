<script lang="ts">
import { IPCEventTypes } from "../util";

import { EventsOn } from "../../wailsjs/runtime/runtime";

import PaneManager from "../components/PaneManager.svelte";
import { userName } from "../store";

EventsOn(IPCEventTypes.Popup, (data) => {
    // TODO
})

</script>

<main>
    <h1>Welcome, {$userName}</h1>
    <PaneManager/>
</main>