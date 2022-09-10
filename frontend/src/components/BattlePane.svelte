<script lang="ts">
import type { PaneInfo } from "../store";
import { currentPaneStore } from "../store";
import BattleScene from "./battle/BattleScene.svelte";
import BattleCommandPalette from "./battle/BattleCommandPalette.svelte";
import BattleChat from "./battle/BattleChat.svelte";

export let info: PaneInfo

$: isFront = info.name === $currentPaneStore

</script>

<main>
    <div hidden={!isFront}>
        <h1>This is a battle pane</h1>
        <div class="left-col">
            <h2>Imagine something cool being here, like a damage calculator</h2>
        </div>
        <div class="center-col">
            <BattleScene/>
            <BattleCommandPalette/>
        </div>
        <div class="right-col">
            <BattleChat roomName={info.name}/>
        </div>
    </div>
</main>