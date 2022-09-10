<script lang="ts">
import type { PaneInfo } from "../store";
import { currentPaneStore } from "../store";
export let info: PaneInfo

$: isFront = info.name === $currentPaneStore
</script>

<main>
    <div hidden={!isFront}>
        <h1>This is a Chat Room Pane</h1>
    </div>
</main>