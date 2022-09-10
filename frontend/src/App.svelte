<script lang="ts">
import { ViewType } from "./util";
import { currentViewStore } from "./store";
import Client from "./views/Client.svelte";
import Login from "./views/Login.svelte";

let currentView: ViewType = ViewType.Login
let lastView: ViewType = currentView

$: {
    lastView = currentView
    currentView = $currentViewStore
}

</script>


<main>
    {#if currentView == ViewType.Login}
    <Login/>
    {:else if currentView == ViewType.Client}
    <Client/>
    <!-- TODO: Settings and Register views -->
    {/if}

</main>