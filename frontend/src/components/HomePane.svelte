<!-- The Home Pane is where the user can search for battles, Users, Rooms, etc-->
<script lang="ts">
import type { PaneInfo } from "../store";
import { currentPaneStore } from "../store"

export let info: PaneInfo

$: isFront = info.name === $currentPaneStore

</script>

<main>
    <div hidden={!isFront}>
        <h1>This is the Home Pane</h1>
        <h2>This page will eventually have something</h2>
    </div>
</main>