<script lang="ts">
import { userName } from "../store";
import { IPCEventTypes, tsPrint, UiEventTypes, ViewType } from "../util";
import { SDLogin } from "../../wailsjs/go/main/App";
import { EventsEmit, EventsOn } from "../../wailsjs/runtime/runtime";

let chooseName: string
let password: string
let errMsg: string = ""

type LoginSuccessEvent = {
    UserName: string
}

type LoginFailEvent = {
    Reason: string
}

function login() {
    tsPrint("sending login")
    SDLogin(chooseName, password)
}

EventsOn(IPCEventTypes.LoginSuccess, (data: LoginSuccessEvent) => {
    tsPrint(`${data.UserName} logged in`)
    userName.update(() => {
        return data.UserName
    })
    EventsEmit(UiEventTypes.ViewChange, ViewType.Client)
})

EventsOn(IPCEventTypes.LoginFail, (data: LoginFailEvent) => {
    tsPrint("login failure")
    errMsg = data.Reason
})

</script>

<main>
    <h1>Enter Username and Password</h1>
    <p id="error-box">{errMsg}</p>
    <input type="text" name="uname" id="uname" bind:value={chooseName}>
    <input type="password" name="pword" id="pword" bind:value={password}>
    <input type="button" value="Login" on:click={login}>
</main>