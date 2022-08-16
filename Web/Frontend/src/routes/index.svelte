<script>
    import { fly } from 'svelte/transition';
    import { onMount } from 'svelte';
    const APIURI = "https://amekanji.com/api/"
    const Pages = {
        Home: 'Home',
        Template: 'Template',
        Input: 'Input',
    };
    let page = Pages.Home;
    let Templates = {
        JPENG : 'JPEng',
        ENGJP : 'EngJP',
        Kanji : 'Kanji',
        Custom : 'Custom'
    }
    let template = Templates.JPENG;

    let customtemplate = ["@{kanjiword} <br> @{kanaword} <br> @{audio} @{CSS}", "@{sense} @{kaniinfoex} @{stroke} @{CSS}"]

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

    let srcinput = "";
    function HandleInput() {
        var lines = srcinput.split('\n');

        var Candidates = []
        var CandidateKanjis = []

        // Fill To Add array
        for(var n = 0; n < lines.length; n++){
            var line = lines[n];
            var segments = line.split(/;|\|/); // Splits on ; or |
            
            var kanji = segments[0];
            kanji = kanji.replace(/\s/g, '');
            
            var kana = segments[1] || "";
            kana = kana.replace(/\s/g, '');
            
            var literal = segments[2] || "";
            literal = literal.replace(/\s/g, '');
            
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

        var Request = {
            "AmeInput" : {
                "Template" : TemplateForm,
                "Input" : InputForm
            }
        }

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
    let Cookie = {};
    function LoadCookies()
    {
        try
        {
            var cookiesrc = document.cookie.substring(8);
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
        document.cookie = "storage=" + cookiesrc + "; expires=Fri, 31 Dec 9999 23:59:59 GMT; SameSite=None;";
    }
    onMount(async () => {
        Cookie = LoadCookies();
        customtemplate = Cookie.ctemplate || ["@{kanjiword} <br> @{kanaword} <br> @{audio} @{CSS}", "@{sense} @{kaniinfoex} @{stroke} @{CSS}"];
        inputarray = Cookie.input || [];
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
        <button type="button" class="btn btn-outline-secondary" on:click={() => {page = Pages.Template}}>
            <i class="bi bi-arrow-down "></i>
        </button>
    </div>
</div>
{:else if page == Pages.Template}
<div id="Template" in:fly={{y: 200, duration: 500, delay: 500}} out:fly={{y:-200, duration: 500}}>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary" on:click={() => {page = Pages.Home}}>
            <i class="bi bi-arrow-up "></i>
        </button>
    </div>
    <div class="content">
        <div class="Content_Title">
            Choose your template!
        </div>
        <div class="template-selector">
            <div class="btn-group" role="group" aria-label="Basic radio toggle button group">
                <input type="radio" class="btn-check" name="btnradio" id="btnradio1" autocomplete="off" bind:group="{template}" value={Templates.JPENG} on:click={() => {template = Templates.JPENG}}>
                <label class="btn btn-outline-secondary" for="btnradio1">JP to EN</label>

                <input type="radio" class="btn-check" name="btnradio" id="btnradio2" autocomplete="off" bind:group="{template}" value={Templates.ENGJP} on:click={() => {template = Templates.ENGJP}}>
                <label class="btn btn-outline-secondary" for="btnradio2">EN to JP</label>

                <input type="radio" class="btn-check" name="btnradio" id="btnradio3" autocomplete="off" bind:group="{template}" value={Templates.Kanji} on:click={() => {template = Templates.Kanji}}>
                <label class="btn btn-outline-secondary" for="btnradio3">Kanji</label>
                
                <input type="radio" class="btn-check" name="btnradio" id="btnradio4" autocomplete="off" bind:group="{template}" value={Templates.Custom} on:click={() => {template = Templates.Custom}}>
                <label class="btn btn-outline-secondary" for="btnradio4">Custom</label>
            </div>
        </div>
        {#if template == Templates.JPENG}
        <div class='card-example'>
            <div class='Kele'><ol><li><div class='Keb'>雨</div></li></ol></div> <style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
            <hr>
            <div class='Rele'><ol><li><div class='Reb'>あめ</div></li></ol></div><br><div class='Sense'><ol><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rain</li></ol></div><div class='example'><div class='lang_jpn'>雨のために彼らは気力をそがれた。</div><div class='lang_eng'>The rain dampened their spirits.</div></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rainy day</li><li>rainy weather</li></ol></div><div class='example'></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>the November suit (in hanafuda)</li></ol></div><div class='example'></div></li></ol></div><br><button type="button" class="btn btn-outline-secondary" style="text-align: center" on:click={() => {playDemo()}}><i class="bi bi-play"></i></button><br><br><div class = 'rexample' id = '0'><div class = 'JP'>明日は雨が降るでしょうか。</div><div class = 'ENG'>Will it rain tomorrow?</div></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
                    </div>
        {:else if template == Templates.ENGJP}
        <div class="card-example">

            <div class='Sense'><ol><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rain</li></ol></div><div class='example'><div class='lang_jpn'>雨のために彼らは気力をそがれた。</div><div class='lang_eng'>The rain dampened their spirits.</div></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rainy day</li><li>rainy weather</li></ol></div><div class='example'></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>the November suit (in hanafuda)</li></ol></div><div class='example'></div></li></ol></div> <style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
            <hr>
            <div class='Kele'><ol><li><div class='Keb'>雨</div></li></ol></div><br><div class='Rele'><ol><li><div class='Reb'>あめ</div></li></ol></div><br><div class='Sense'><ol><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rain</li></ol></div><div class='example'><div class='lang_jpn'>雨のために彼らは気力をそがれた。</div><div class='lang_eng'>The rain dampened their spirits.</div></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>rainy day</li><li>rainy weather</li></ol></div><div class='example'></div></li><li><div class='pos'><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class='gloss'><ol><li>the November suit (in hanafuda)</li></ol></div><div class='example'></div></li></ol></div><br><br><button type="button" class="btn btn-outline-secondary" style="text-align: center" on:click={() => {playDemo()}}><i class="bi bi-play"></i></button><div class = 'rexample' id = '0'><div class = 'JP'>明日は雨が降るでしょうか。</div><div class = 'ENG'>Will it rain tomorrow?</div></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
        </div>
        {:else if template == Templates.Kanji}
        <div class='card-example'>


            <div class = 'kliteral'>雨</div> <style>.kanji_info{text-align:center;list-style-position:inside}.kanji_info>ol>li{color:#333;font-size:100%}.kanji_info>.kanji_instance>.literal{color:#000;font-size:150%;text-align:center}.kanji_info>.kanji_instance>.meanings{text-align:center;margin-right:3%}.kanji_info>.kanji_instance>.meanings>ol>li{font-size:100%}.kanji_info>.kanji_instance>.meanings>ol>li:nth-child(1){font-size:125%}.kanji_info>.kanji_instance>.readings>ul>li{display:inline;margin-left:4%;margin-right:4%}.kanji_info>.kanji_instance>.misc>*{margin-top:3%;margin-bottom:3%}.kliteral{text-align: center;color:#000;font-size:150%}</style>
            <br>
            <div class = 'kanji_info'><div class = 'kanji_instance'><div class = 'literal'>雨</div><div class = meanings><ol><li>rain</li></ol></div><div class = readings><ul><li>ウ</li><li>あめ</li><li>あま-</li><li>-さめ</li></ul></div><div class='misc'><div class='grade'> Grade: 1</div><div class='strokecount'> Stroke count: 8</div><div class='jlpt'> JLPT: 4</div><div class='freq'> Frequency: 950</div></div></div></div><br><div class = 'stroke_set'><div class = 'stroke ANDAS'><img alt="stroke" src='./media/ANDAS3561.gif'></div></div><style>.kanji_info{text-align:center;list-style-position:inside}.kanji_info>ol>li{color:#333;font-size:100%}.kanji_info>.kanji_instance>.literal{color:#000;font-size:150%;text-align:center}.kanji_info>.kanji_instance>.meanings{text-align:center;margin-right:3%}.kanji_info>.kanji_instance>.meanings>ol>li{font-size:100%}.kanji_info>.kanji_instance>.meanings>ol>li:nth-child(1){font-size:125%}.kanji_info>.kanji_instance>.readings>ul>li{display:inline;margin-left:4%;margin-right:4%}.kanji_info>.kanji_instance>.misc>*{margin-top:3%;margin-bottom:3%}.kliteral{text-align: center;color:#000;font-size:150%}</style>
        </div>
        {:else}
        <div class='template-creator'>
            {#each customtemplate as form, fid}
                <div class="mb-3">
                    <label for="form_{fid}" class="form-label">Field {fid + 1}: </label>
                    <textarea class="form-control" id="form_{fid}" rows="3" bind:value={customtemplate[fid]} on:change={() => {UpdateCookie();}}></textarea>
                    {#if customtemplate.length > 2}
                        <button type="button" class="btn btn-outline-secondary" on:click={() => {customtemplate.splice(fid, 1); customtemplate = customtemplate; UpdateCookie();}}> <i class="bi bi-eraser"></i> </button>
                    {/if}
                </div>
            {/each}
            <button type="button" class="btn btn-outline-secondary" on:click={() => {customtemplate.push(""); customtemplate = customtemplate;UpdateCookie();}}> <i class="bi bi-plus"></i> </button>
        </div>
        {/if}
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
        <button type="button" class="btn btn-outline-secondary" on:click={() => {page = Pages.Template}}>
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
