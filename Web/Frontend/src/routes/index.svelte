<script>
    import { fly } from 'svelte/transition';
    const Pages = {
        Home: 'Home',
        Template: 'Template',
        Input: 'Input',
    };
    let page = Pages.Home;
    const Templates = {
        JPENG : 'JPEng',
        ENGJP : 'EngJP',
        Kanji : 'Kanji',
        Custom : 'Custom'
    }
    let template = Templates.JPENG;

    let customtemplate = ["@{kanjiword} <br> @{kanaword} <br> @{audio} @{css}", "@{sense} @{kaniinfoex} @{stroke} @{css}"]

    function playDemo() {
        let audioFile = new Audio("./media/sangyou.mp3");;
        audioFile.play()
    }
    function SendForm() {

    }
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

{#if page == Pages.Home}
<div id="Home" in:fly={{y: 200, duration: 500, delay: 500}} out:fly={{y:-200, duration: 500}}>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary disabled" on:click={() => {page = Pages.Home}}>
            <i class="bi bi-arrow-up "></i>
        </button>
    </div>
    <div class="content">
        Home
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
                <input type="radio" class="btn-check" name="btnradio" id="btnradio1" autocomplete="off" checked on:click={() => {template = Templates.JPENG}}>
                <label class="btn btn-outline-secondary" for="btnradio1">JP to EN</label>

                <input type="radio" class="btn-check" name="btnradio" id="btnradio2" autocomplete="off" on:click={() => {template = Templates.ENGJP}}>
                <label class="btn btn-outline-secondary" for="btnradio2">EN to JP</label>

                <input type="radio" class="btn-check" name="btnradio" id="btnradio3" autocomplete="off" on:click={() => {template = Templates.Kanji}}>
                <label class="btn btn-outline-secondary" for="btnradio3">Kanji</label>
                
                <input type="radio" class="btn-check" name="btnradio" id="btnradio4" autocomplete="off" on:click={() => {template = Templates.Custom}}>
                <label class="btn btn-outline-secondary" for="btnradio4">Custom</label>
            </div>
        </div>
        {#if template == Templates.JPENG}
        <div class='card-example'>
            <div class="Kele"><ol><li><div class="Keb">産業</div></li></ol></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
            <hr>
            <div class="Rele"><ol><li><div class="Reb">さんぎょう</div></li></ol></div><br><div class="Sense"><ol><li><div class="pos"><ul><li>noun (common) (futsuumeishi)</li><li>noun (common) (futsuumeishi)ounoun (common) (futsuumeishi)s which may take the genoun (common) (futsuumeishi)itive case particle 'noun (common) (futsuumeishi)o'</li></ul></div><div class="gloss"><ol><li>industry</li></ol></div><div class="example"><div class="lang_jpn">日本は産業国である。</div><div class="lang_eng">Japan is an industrial country.</div></div></li><li><div class="pos"><ul><li>noun (common) (futsuumeishi)</li></ul></div><div class="gloss"><ol><li>livelihood</li><li>occupation</li></ol></div><div class="example"></div></li></ol></div><br><div class="audio" style='text-align: center'><button type="button" class="btn btn-outline-secondary" style="text-align: center" on:click={() => {playDemo()}}><i class="bi bi-play"></i></button></div><br><br><div class="rexample" id="0"><div class="JP">日本は産業国である。</div><div class="ENG">Japan is an industrial country.</div></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style></div>
        {:else if template == Templates.ENGJP}
        <div class="card-example">
            <div class="Sense"><ol><li><div class="pos"><ul><li>noun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)</li><li>noun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)ounoun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)s which may take the genoun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)itive case particle 'noun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)o'</li></ul></div><div class="gloss"><ol><li>industry</li></ol></div><div class="example"><div class="lang_jpn">日本は産業国である。</div><div class="lang_eng">Japan is an industrial country.</div></div></li><li><div class="pos"><ul><li>noun (common) (futsuumeishi)ounoun (common) (futsuumeishi) (commonoun (common) (futsuumeishi)) (futsuumeishi)</li></ul></div><div class="gloss"><ol><li>livelihood</li><li>occupation</li></ol></div><div class="example"></div></li></ol></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
            <hr>
            <div class="Kele"><ol><li><div class="Keb">産業</div></li></ol></div><br><div class="Rele"><ol><li><div class="Reb">さんぎょう</div></li></ol></div><br><button type="button" class="btn btn-outline-secondary" style="text-align: center" on:click={() => {playDemo()}}><i class="bi bi-play"></i></button><br><br><div class="rexample" id="0"><div class="JP">日本は産業国である。</div><div class="ENG">Japan is an industrial country.</div></div><style>.Kele{text-align:center;list-style-position:inside}.Kele>ol>li{color:#333;font-size:100%}.Kele>ol>li:nth-child(1){color:#000;font-size:150%}.Kele>ol>li:nth-child(2){color:#131313;font-size:125%}.Keb{display:inline-block;text-align:left}.Rele{text-align:center;list-style-position:inside}.Rele>ol>li{color:#333;font-size:100%}.Rele>ol>li:nth-child(1){color:#000;font-size:150%}.Rele>ol>li:nth-child(2){color:#131313;font-size:125%}.Reb{display:inline-block;text-align:left}.Sense>ol>li{text-align:left;margin-left:10%}.pos>ul{list-style-type:disc}.pos>ul>li{color:gray;font-size:75%;display:inline-flex;margin-left:1em;margin-right:1em}.gloss>ol>li{margin-top:1em;margin-bottom:1em}.gloss>ol>li:nth-child(1){font-size:115%}.example>.lang_jpn{font-size:100%;color:#1d1d29;margin-top:1em}.example>.lang_eng{font-size:75%;color:#303030;margin-bottom:1em}.example{margin-top:1em;margin-bottom:1em}</style><style>.rexample{text-align:left;margin-left:10%}.rexample>.JP{font-size:100%;color:#1d1d29;margin-top:1em}.rexample>.ENG{font-size:75%;color:#303030;margin-bottom:1em}</style>
        </div>
        {:else if template == Templates.Kanji}
        <div class='card-example'>
            <div class="kliteral">絹</div><style>.kanji_info{text-align:center;list-style-position:inside}.kanji_info>ol>li{color:#333;font-size:100%}.kanji_info>.kanji_instance>.literal{color:#000;font-size:150%;text-align:center}.kanji_info>.kanji_instance>.meanings{text-align:center;margin-right:3%}.kanji_info>.kanji_instance>.meanings>ol>li{font-size:100%}.kanji_info>.kanji_instance>.meanings>ol>li:nth-child(1){font-size:125%}.kanji_info>.kanji_instance>.readings>ul>li{display:inline;margin-left:4%;margin-right:4%}.kanji_info>.kanji_instance>.misc>*{margin-top:3%;margin-bottom:3%}.kliteral{text-align: center;color:#000;font-size:150%}</style>
            <hr>
            <div class="kanji_info"><div class="kanji_instance"><div class="literal">絹</div><div class="meanings"><ol><li>silk</li></ol></div><div class="readings"><ul><li>ケン</li><li>きぬ</li></ul></div><div class="misc"><div class="grade"> Grade: 6</div><div class="strokecount"> Stroke count: 13</div><div class="jlpt"> JLPT: 1</div><div class="freq"> Frequency: 1916</div></div></div></div><br><div class="stroke_set"><div class="stroke ANDAS"><img src="./media/ANDAS1361.gif"></div></div><style>.kanji_info{text-align:center;list-style-position:inside}.kanji_info>ol>li{color:#333;font-size:100%}.kanji_info>.kanji_instance>.literal{color:#000;font-size:150%;text-align:center}.kanji_info>.kanji_instance>.meanings{text-align:center;margin-right:3%}.kanji_info>.kanji_instance>.meanings>ol>li{font-size:100%}.kanji_info>.kanji_instance>.meanings>ol>li:nth-child(1){font-size:125%}.kanji_info>.kanji_instance>.readings>ul>li{display:inline;margin-left:4%;margin-right:4%}.kanji_info>.kanji_instance>.misc>*{margin-top:3%;margin-bottom:3%}.kliteral{text-align: center;color:#000;font-size:150%}</style>
        </div>
        {:else}
        <div class='template-creator'>
            {#each customtemplate as form, fid}
                <div class="mb-3">
                    <label for="form_{fid}" class="form-label">Field {fid + 1}: </label>
                    <textarea class="form-control" id="form_{fid}" rows="3" bind:value={customtemplate[fid]}></textarea>
                    {#if customtemplate.length > 2}
                        <button type="button" class="btn btn-outline-secondary" on:click={() => {customtemplate.splice(fid, 1); customtemplate = customtemplate;}}> <i class="bi bi-eraser"></i> </button>
                    {/if}
                </div>
            {/each}
            <button type="button" class="btn btn-outline-secondary" on:click={() => {customtemplate.push(""); customtemplate = customtemplate;}}> <i class="bi bi-plus"></i> </button>
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
        Input
    </div>
    <div class='controller'>
        <button type="button" class="btn btn-outline-secondary" on:click={SendForm}>
            <i class="bi bi-arrow-right "></i>
        </button>
    </div>
</div>
{/if}

<style>
    .controller {
        text-align: center;
        margin-top: 30px;
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
</style>