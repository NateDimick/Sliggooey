<script lang="ts">
import { IPCEventTypes } from "../util";

import { EventsOn } from "../../wailsjs/runtime/runtime";

import PaneManager from "../components/PaneManager.svelte";
import { userName } from "../store";


let user: string
userName.subscribe((name: string) => {
    user = name
})

EventsOn(IPCEventTypes.Popup, (data) => {
    // TODO
})

</script>

<main>
    <h1>Welcome, {user}</h1>
    <PaneManager/>
</main>