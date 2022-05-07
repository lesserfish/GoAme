<script>
    import { page } from '$app/stores';
    const uuid = $page.url.searchParams.get('id');
    
    function poll() {
    
    }

    var message = "Sending request"
    var status = "loading";
</script>

<div id="header">
    <nav class="navbar navbar-light bg-light">
        <div class="container-fluid">
            <a class="navbar-brand" href="./">
                <img src="./logo_64.png" alt="" width="30" height="24" class="d-inline-block align-text-top">
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

