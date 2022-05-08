<script>
    import { onMount } from 'svelte';

    const APIURI = "http://localhost:9000/"
    const timeoutdelay = 2000; // 2 seconds
    var message = "Sending request"
    var status = "loading";
    
    function poll(uuid) {
        
        var uri = APIURI + "get?id=" + uuid;
        var xmlHttpRequest = new XMLHttpRequest();

        xmlHttpRequest.open("GET", uri, true);
        xmlHttpRequest.onload = function(){
            var statuscode = this.status;
            var rawresponse = this.responseText;
            
            var uri = APIURI + "get?id=" + uuid;
            
            var response = {};
            var contentType = this.getResponseHeader('Content-Type');
            if(contentType == "application/zip") {
                status = "success";
                message = "";
                window.open(uri, '_self');
                return;
            }
            else if(contentType == "application/json"){
                response = JSON.parse(rawresponse);
            }
            
            if(statuscode != 200){
                status = "error";
                message = "Error:  " + (response.Message || "Invalid response from server.");
                return;
            }

            var taskstatus = response.Status;
            var progress = Math.round(response.Progress * 100);

            if(taskstatus == "Deleted"){
                status = "error"
                message = "Error. File has already been deleted."
                return;
            }
            else if(taskstatus == "In Progress") {
                status = "loading";
                message = "In progress.... (" + String(progress) + "%)";
                setTimeout(() => {poll(uuid)}, timeoutdelay);
                return;
            }
            else if(taskstatus == "Failed") {
                status = "error";
                message = "Task failed :(";
                return;
            }
            else if(taskstatus == "Accepted") {
                status = "loading";
                message = "In progress...";
                setTimeout(() => {poll(uuid)}, timeoutdelay);
                return;
            } else {
                console.log(taskstatus);
                console.log(progress);
            }

        }

        xmlHttpRequest.send();
    }
    onMount(async () => {
        const params = new URLSearchParams(window.location.search);
        const uuid = params.get("id");
        poll(uuid);
	});

</script>

<div id="header">
    <nav class="navbar navbar-light bg-light">
        <div class="container-fluid">
            <a class="navbar-brand" href="./">
                <img src="./logo_64.png" alt="logo 64px" width="30" height="24" class="d-inline-block align-text-top">
                AmeKanji
            </a>
        </div>
    </nav>
</div>

<div class="container">
    <div class="row loading">
        {#if status=="loading"}
            <div class="spinner-grow text-primary" role="status">
                <span class="sr-only"></span>
            </div>
        {:else if status=="error"}
            <i class="bi bi-x-circle-fill text-danger" style="scale: 200%;"></i>
        {:else if status=="success"}
            <i class="bi bi-suit-heart-fill text-danger" style="scale: 200%;"></i>
        {/if}
    </div>
    <div class="row message">
        <div class="text">
            <p>{message}</p>
        </div>
    </div>
</div>


<style>
    .container {
        width: 100%;
        height: 100%;
        margin: 0;
        left: 0;
        right: 0;
    }
    .loading {
        position: absolute;
        left: 50%;
        top: 50%;
    }
    .message {
        position: absolute;
        top: 60%;
        margin: 0;
        left: 0;
        width: 100%;
    }
    .text {
        text-align: center;
    }
</style>

