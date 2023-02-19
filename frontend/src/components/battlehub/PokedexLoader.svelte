<script lang="ts">
import { tsPrint } from "../../util";
import { pokedex, specialIconIds } from "../../store"
import { GetPokedex, GetSpecialSpriteNumbers } from "../../wailsjs/go/backend/App";

async function loadPokedex() {
    if (Object.keys($pokedex).length === 0) {
        tsPrint("fetching pokedex from showdown server")
        let resp = await GetPokedex()
        let dex = JSON.parse(resp)
        tsPrint(`got dex from showdown: ${JSON.stringify(dex.bulbasaur)}`)
        pokedex.set(dex)
        tsPrint("pokedex loading complete")
        let resp2 = await GetSpecialSpriteNumbers()
        specialIconIds.set(resp2)
        tsPrint("special pokedex icon id loading complete")
    }
}

loadPokedex()
</script>