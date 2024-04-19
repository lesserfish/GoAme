<script>
    import { fly } from 'svelte/transition';
    import { onMount } from 'svelte';
    const APIURI = window.location.origin + "/api/";
    const Pages = {
        Home: 'Home',
        Input: 'Input',
    };
    let page = Pages.Home;

    let inputarray = []
    let tag = "AmeKanji"
    let allSelected = false;
    
    //var id = 0;
    function NewID(){
        var id = 0;
        for(var i = 0; i < inputarray.length; i++)
        {
            id = id > inputarray[i].id ? id : inputarray[i].id + 1;
        }
        return id;
    }
    function SelectionChange() {
        allSelected = true;
        for(var i = 0; i < inputarray.length; i++) {
            if(inputarray[i].selected == false) {
                allSelected = false;
                return;
            }
        }
    }
    function SelectAll() {
        for(var i = 0; i < inputarray.length; i++) {
            inputarray[i].selected = allSelected;
        }
    }
    function DeleteSelected() {
        for(var i = 0; i < inputarray.length; i++) {
            if(inputarray[i].selected == true){
                inputarray.splice(i, 1);
                inputarray = inputarray;
                return DeleteSelected() + 1;
            }
        }
        UpdateCookie();
        return 0;
    }
    function playDemo() {
        let audioFile = new Audio("./media/ame.mp3");;
        audioFile.play()
    }

    const allkana = ['ー', 'ぁ', 'あ', 'ぃ', 'い', 'ぅ', 'う', 'ぇ', 'え', 'ぉ', 'お', 'か', 'が', 'き', 'ぎ', 'く', 'ぐ', 'け', 'げ', 'こ', 'ご', 'さ', 'ざ', 'し', 'じ', 'す', 'ず', 'せ', 'ぜ', 'そ', 'ぞ', 'た', 'だ', 'ち', 'ぢ', 'っ', 'つ', 'づ', 'て', 'で', 'と', 'ど', 'な', 'に', 'ぬ', 'ね', 'の', 'は', 'ば', 'ぱ', 'ひ', 'び', 'ぴ', 'ふ', 'ぶ', 'ぷ', 'へ', 'べ', 'ぺ', 'ほ', 'ぼ', 'ぽ', 'ま', 'み', 'む', 'め', 'も', 'ゃ', 'や', 'ゅ', 'ゆ', 'ょ', 'よ', 'ら', 'り', 'る', 'れ', 'ろ', 'ゎ', 'わ', 'ゐ', 'ゑ', 'を', 'ん', 'ァ', 'ア', 'ィ', 'イ', 'ゥ', 'ウ', 'ェ', 'エ', 'ォ', 'オ', 'カ', 'ガ', 'キ', 'ギ', 'ク', 'グ', 'ケ', 'ゲ', 'コ', 'ゴ', 'サ', 'ザ', 'シ', 'ジ', 'ス', 'ズ', 'セ', 'ゼ', 'ソ', 'ゾ', 'タ', 'ダ', 'チ', 'ヂ', 'ッ', 'ツ', 'ヅ', 'テ', 'デ', 'ト', 'ド', 'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ', 'バ', 'パ', 'ヒ', 'ビ', 'ピ', 'フ', 'ブ', 'プ', 'ヘ', 'ベ', 'ペ', 'ホ', 'ボ', 'ポ', 'マ', 'ミ', 'ム', 'メ', 'モ', 'ャ', 'ヤ', 'ュ', 'ユ', 'ョ', 'ヨ', 'ラ', 'リ', 'ル', 'レ', 'ロ', 'ヮ', 'ワ', 'ヰ', 'ヱ', 'ヲ', 'ン', 'ヴ', 'ヵ', 'ヶ', '゛', '゜']
    
    function is_kanji(chr) {
        if(chr.charCodeAt(0) >= 0x4E00 && chr.charCodeAt(0) <= 0x9fff) {
            return true
        }
        return false
    }

    function clean_kanji(input) {
        var valid = false
        for(var i = 0; i < input.length; i++)
        {
            var chr = input[i]
            if(is_kanji(chr)) {
                valid = true
                break
            }
        }
        if(valid) {
            return input
        }

        return ""
    }
    let srcinput = "";
    function HandleInput() {
        var lines = srcinput.split('\n');

        var Candidates = []
        var CandidateKanjis = []

        // Fill To Add array
        for(var n = 0; n < lines.length; n++){
            var line = lines[n];
            var segments = line.split(/;|\||:/); // Splits on ; or |
            

            var kanji = "";
            var kana = "";
            var literal = "";
            if(segments.length == 1)
            {
                // Check if there are kanji in the word. If there aren't, just fill kana
                
                var word = segments[0];
                var ver = clean_kanji(word)
                if(ver == "") {
                    kanji = ""
                    kana = word
                } else {
                    kanji = word
                    kana = ""
                }
            }
            else {
                kanji = segments[0];
                kanji = kanji.replace(/\s/g, '');
                kanji = kanji.replace(/ /g, '');
                kanji = clean_kanji(kanji)
                
                kana = segments[1] || "";
                kana = kana.replace(/\s/g, '');
                kana = kana.replace(/ /g, '');
                
                literal = segments[2] || "";
                literal = literal.replace(/\s/g, '');
                literal = literal.replace(/ /g, '');
            }

            kanji = kanji.replace("<br>", "");
            kana = kana.replace("<br>", "");
            literal = literal.replace("<br>", "");
            
            if(kanji == "" && kana == "" && literal == ""){
                continue;
            }

            if(literal == "") {
                for(var k = 0; k < kanji.length; k++){
                    var letter = kanji[k];
                    var hexval = letter.codePointAt(0);
                    if(hexval >= 0x4E00 && hexval <= 0x9FAF) {
                        literal += letter;
                    }
                }
            }

            Candidates.push({ 'kanji': kanji,
                        'kana': kana,
                        'literal': literal});
            CandidateKanjis.push(kanji);
        }
        // Download helpful info from rest api
        
        var uri = APIURI + "help"

        var xmlHttpRequest = new XMLHttpRequest();
        xmlHttpRequest.open('POST', uri, true);
        xmlHttpRequest.setRequestHeader('Content-Type', 'application/json')
        xmlHttpRequest.onload = function(){
            var rawresponse = this.responseText;
            var response = JSON.parse(rawresponse);

            if(this.status != 200){
                console.error(response.Message);
                return;
            }
            var info = response.Response;

            for(var k = 0; k < Candidates.length; k++) {
                var currentCandidate = Candidates[k];
                var newEntry = { 
                    id : NewID(), 
                    selected: false,
                    kanji: currentCandidate.kanji,
                    kana: currentCandidate.kana,
                    literal: currentCandidate.literal,
                    kanadb: info[currentCandidate.kanji] || []
                }

                inputarray.push(newEntry);
            }
            inputarray = inputarray;
            UpdateCookie();
            setTimeout(() => {SetLoading('hide')}, 600);
        }
        var requestbody = JSON.stringify(CandidateKanjis);

        SetLoading('show');
        xmlHttpRequest.send(requestbody);
        srcinput = "";
    }
    function HandleInputKey(keyevent) {
        if(keyevent.key == 'Enter') {
            keyevent.preventDefault();
            HandleInput();
        }
    }
    function HandleClipboard(clipevent){
        clipevent.preventDefault();
        var data = (clipevent.originalEvent || clipevent).clipboardData.getData('text/plain');
        srcinput = data;
        HandleInput();
    }
    function SetLoading(option) {
        window.$('#loadingModal').modal(option);
    }
    function SendForm() {
        var chosentemplate = []
        
        if(template == Templates.JPENG){
                chosentemplate = ["@{kanjiword} @{CSS}","@{kanaword}<br>@{sense}<br>@{audio}<br><br>@{example}@{CSS}"]
        }else if(template ==  Templates.ENGJP){
                chosentemplate = ["@{sense} @{CSS}","@{kanjiword}<br>@{kanaword}<br>@{sense}<br><br>@{audio}@{example}@{CSS}"]
        }else if(template == Templates.Kanji){
                chosentemplate = ["@{literal} @{CSS}","@{kanjiinfoex}<br>@{stroke}@{CSS}"]
        } else if(template ==  Templates.Custom){
                chosentemplate = customtemplate;
        }
        
        var TemplateForm = {"Fields" : chosentemplate, "Tag" : tag};
        var InputForm = [];

        for(var i = 0; i < inputarray.length; i++) {
            var currentinput = inputarray[i]
            var chosenkana = currentinput.kana;
            if(chosenkana == ""){
                chosenkana = currentinput.kanadb[0] || "";
            }
            var entry = {
                "kanjiword" : currentinput.kanji,
                "kanaword" : chosenkana,
                "literal" : currentinput.literal
            }

            InputForm.push(entry);
        }

        var Request = InputForm;

        var uri = APIURI + "post"

        var xmlHttpRequest = new XMLHttpRequest();
        xmlHttpRequest.open('POST', uri, true);
        xmlHttpRequest.setRequestHeader('Content-Type', 'application/json')
        xmlHttpRequest.onload = function(){
            var rawresponse = this.responseText;

            var response = JSON.parse(rawresponse);
             
            if(this.status != 200){
                console.error(response.Message);
                // TODO: Handle Error!
                return;
            }

            var uuid = response.UUID;

            // Redirect to uuid

            var redirectionuri = "get.html?id=" + uuid;
            window.location.replace(redirectionuri);
        }
        
        var requestbody = JSON.stringify(Request);

        xmlHttpRequest.send(requestbody);

    }

    function getCookie(name) {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop().split(';').shift();
    }
    let Cookie = {};
    function LoadCookies()
    {
        try
        {
            //var cookiesrc = document.cookie.substring(8);
            //var cookiesrc = getCookie("storage");
            var cookiesrc = localStorage.getItem("storage") || JSON.stringify({});
            var obj = JSON.parse(cookiesrc);
            return obj;
        } catch (e)
        {
            console.error(e);
            return {}
        }
    }
    function SetCookie()
    {
        var cookiesrc = JSON.stringify(Cookie);
        //document.cookie = "storage=" + cookiesrc + "; expires=Fri, 31 Dec 9999 23:59:59 GMT; SameSite=None; Secure;";
        localStorage.setItem("storage", cookiesrc);
    }
    onMount(async () => {
        Cookie = LoadCookies();
        customtemplate = Cookie.ctemplate || ["@{kanjiword} <br> @{kanaword} <br> @{audio} @{CSS}", "@{sense} @{kaniinfoex} @{stroke} @{CSS}"];
        inputarray = Cookie.input || [];

        // Check that no duplicate ids exist.
        for(var i = 0; i < inputarray.length; i++)
        {
            inputarray[i].id = i;
        }
    });
    function UpdateCookie()
    {
        Cookie.input = inputarray;
        Cookie.ctemplate = customtemplate;
        SetCookie();
    }
    
</script>

<div id="header">
    <nav class="navbar navbar-light bg-light">
        <div class="container-fluid">
            <a class="navbar-brand" href="./">
                <img alt="logo" src="./logo_64.png" width="30" height="24" class="d-inline-block align-text-top">
                AmeKanji
            </a>
        </div>
    </nav>
</div>

{#if page == Pages.Home}
<div id="Home" in:fly={{y: 200, duration: 500, delay: 500}} out:fly={{y:-200, duration: 500}}>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary disabled" on:click={() => {page = Pages.Home}}>
            <i class="bi bi-arrow-up "></i>
        </button>
    </div>
    <div class="content">
        <div class="home">
            <div class="logo">
                <h1>雨</h1>
            </div>
            <div class="description">
                AmeKanji! An anki tool deck creator.
            </div>
        </div>
    </div>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary" on:click={() => {page = Pages.Input}}>
            <i class="bi bi-arrow-down "></i>
        </button>
    </div>
</div>
{:else if page == Pages.Input}
<div id="Input" in:fly={{y: 200, duration: 500, delay: 500}} out:fly={{y:-200, duration: 500}}>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary" on:click={() => {page = Pages.Home}}>
            <i class="bi bi-arrow-up "></i>
        </button>
    </div>
    <div class="content">
        <div class="row justify-content-center" style="margin-bottom: 20px;">
           <div class='col-6 align-self-center'>
               <div class="cell description edit inputbox" contenteditable="true" bind:innerHTML="{srcinput}" on:input="{() => {}}" on:paste="{HandleClipboard}" on:keydown="{HandleInputKey}" placeholder="食べる|たべる"></div>
           </div>
        </div>
        <div class="row input_field">
            <div class="entry_container">
                <div class="row allselect">
                    <div class="col allselectleft">
                        <input id="selectallcheck" class="form-check-input" type="checkbox" value="" bind:checked="{allSelected}" on:change="{SelectAll}">
                        <label for="selectallcheck">Select all</label>
                    </div>
                    <div class="col allselectright">
                        <button type="button" class="btn btn-sm btn-outline-danger" on:click="{DeleteSelected}"><i class="bi bi-x">Erase selected</i></button>
                    </div>
                </div>
                {#each inputarray as entry (entry.id)}
                    <div class="entry">
                        <input class="form-check-input" type="checkbox" value="" bind:checked="{entry.selected}" on:change="{() => {SelectionChange();}}">
                        <input disabled type="text" bind:value="{entry.kanji}" placeholder="kanji reading">
                        <input type="text" bind:value="{entry.kana}" placeholder="{entry.kanadb[0] || 'kana reading'}"
                            list="entry_{entry.id}_candidates">
                        <input type="text" bind:value="{entry.literal}">
                        <button type="button" class="btn btn-sm btn-outline-danger" on:click={() => {inputarray.splice(inputarray.indexOf(entry), 1); inputarray = inputarray; UpdateCookie();}}><i class="bi bi-x"></i></button>
                        {#if entry.kanadb.length > 0}
                            <datalist id="entry_{entry.id}_candidates">
                            {#each entry.kanadb as reading, rid}
                                <option value={reading}>{reading}</option>
                            {/each}
                            </datalist>
                        {/if}
                    </div>
                {/each}
            </div>
        </div>
        <div class="tagfield">
             <label for="tag">Tag: </label>
             <input id="tag" type="text" bind:value="{tag}" placeholder="Tag for your cards">
        </div>
    </div>
    <div class='controller'>
        <button type="button" class="btn btn-outline-primary" on:click={SendForm}>
            <i class="bi bi-arrow-right "></i>
        </button>
    </div>
    <div class="loadingpopup">
        <div class="modal fade" id="loadingModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
            <div class="modal-dialog" role="document">
                <div class="modal-content">
                    <div class="d-flex justify-content-center">
                        <div class="spinner-border text-primary" role="status">
                            <span class="sr-only"></span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{/if}

<style>
    .controller {
        text-align: center;
        margin-top: 30px;
        margin-bottom: 15px;
    }
    .content {
        text-align: center;
    }
    .card-example {
        left: 50%;
        margin: auto;
        height: 18em;
        width: 80%; 
        overflow-x: hidden;
        overflow-y: auto;
        text-align: center;
        outline: 3px solid rgb(209, 209, 209);
        box-shadow: 0 10px rgb(221, 221, 221);
    }
    .template-creator {
        left: 50%;
        margin: auto;
        height: 18em;
        width: 80%; 
        overflow-x: hidden;
        overflow-y: auto;
        text-align: center;
        outline: 3px solid rgb(209, 209, 209);
        box-shadow: 0 10px rgb(221, 221, 221);
    }
    .allselectleft {
        text-align: left;
        margin-left: 25%;
    }
    .allselectright > button {
        text-align: right;
        margin-right: 35%;
    }
    .allselect {
        margin-bottom: 15px;
    }
    .inputbox {
        border: 1px solid black;
        text-align: left;
    }
    .loadingpopup {
        text-align: center;
    }
    .loadingpopup * {
        outline: 0px solid black !important;
    }
    .modal-content {
        background-color: transparent;
        outline: 0;
        border: 0;
        box-shadow: 0 0 0 0;
    }
    .logo h1 {
        font-size: 168px;
        font-family: 'Yuji Boku', serif;
    }

</style>
