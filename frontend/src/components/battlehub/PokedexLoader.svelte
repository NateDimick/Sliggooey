<script lang="ts">
import { tsPrint } from "../../util";
import { pokedex } from "../../store"
import { GetPokedex } from "../../wailsjs/go/backend/App";

async function loadPokedex() {
    if (Object.keys($pokedex).length === 0) {
        tsPrint("fetching pokedex from showdown server")
        let resp = await GetPokedex()
        let dex = JSON.parse(resp)
        tsPrint(`got dex from showdown: ${JSON.stringify(dex.bulbasaur)}`)
        pokedex.set(dex)
        tsPrint("pokedex loading complete")
    }
}

loadPokedex()
</script>