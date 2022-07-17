<script lang="ts">
import { EventsOn } from "../wailsjs/runtime/runtime";
import { UiEventTypes, ViewType } from "./util";
import Client from "./views/Client.svelte";
import Login from "./views/Login.svelte";

let currentView: ViewType = ViewType.Login
let lastView: ViewType = currentView

EventsOn(UiEventTypes.ViewChange, (data: ViewType) => {
    if (data === ViewType.Previous) {
        let tmp: ViewType = currentView
        currentView = lastView
        lastView = tmp
    } else {
        lastView = currentView
        currentView = data
    }
})

</script>


<main>
    {#if currentView == ViewType.Login}
    <Login/>
    {:else if currentView == ViewType.Client}
    <Client/>
    <!-- TODO: Settings and Register views -->
    {/if}

</main>