<!-- Pane Tab is a tab that allows the user to selct the specific pane they want to see by clicking this -->
<script lang="ts">
    import { UiEventTypes, PaneType } from "../util"
    import { EventsEmit } from "../../wailsjs/runtime/runtime"
    import { currentPaneStore, PaneInfo } from "../store"

    export let info: PaneInfo

    function seeThisPane() {
        currentPaneStore.set(info.name)
    }

    function deleteThisPane() {
        EventsEmit(UiEventTypes.DeletePane, info.name)
    }
</script>

<main>
    <input type="button" value={info.name} on:click={seeThisPane}>
    {#if info.type !== PaneType.HomePane && info.type !== PaneType.ChatPane}
        <input type="button" value="&#10006" on:click={deleteThisPane}>
    {/if}
</main>