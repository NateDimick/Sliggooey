<script lang="ts">
import type { PaneInfo } from "../store";
import { currentPaneStore } from "../store";
export let info: PaneInfo

$: isFront = info.name === $currentPaneStore
</script>

<main>
    <div hidden={!isFront}>
        <h1>Chat Rooms</h1>
        <p>Imagine eventually there is a list of public chat rooms listed here that are also joinable</p>
    </div>
</main>