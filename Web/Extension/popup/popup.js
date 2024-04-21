class Popup {
    async Main(){
        var storage = new Storage();
        this.storage = storage;
        var memory = await storage.LoadMemory();
        this.memory = memory;

        this.#Go();
    }
    
    // Reloads the memory
    async #ReloadMemory()
    {
        var memory = await this.storage.LoadMemory();
        this.memory = memory;
    }

    // Do everything
    async #Go() {
        this.#Initialize();
        this.#PopulateEntries();
    }

    // Initialized the required variables
    async #Initialize(){
        this.prevCheckboxID = 0;
        this.prevChecked = true;
        this.#InitializeTrash();
        this.#InitializeMenu();
        this.#InitializeInput();
        this.#InitializeDownload();
    }
    // Loads data from browser and load entries into the HTML
    async #PopulateEntries(){
        for(var id = 0; id < this.memory.length; id++){
            var entry_object = this.memory[id];

            var entryDiv = $('<div>').addClass('entry');
            entryDiv.attr('data-id', id);

            var checkbox = $('<input>').addClass('entry-checkbox').attr('type', 'checkbox');
            checkbox.on('click', (e) => this.#HandleCheckbox(e));

            var wordDiv = $('<div>').addClass('word').text(entry_object.w);
            var readingDiv = $('<div>').addClass('reading').text(entry_object.k); 
            

            entryDiv.append(checkbox, wordDiv, readingDiv);
            $('.entries').append(entryDiv);
        }
    }
    
    // Everytime a checkbox is clicked, this function is evoked
    async #HandleCheckbox(e){
        var checkbox = $(e.target);
        var checked = checkbox.prop("checked");
        var id = checkbox.parent().attr("data-id");

        if(e.shiftKey && id != this.prevCheckboxID){
            var that = this;

            $(".entry").each(function(){

                var this_id = $(this).attr("data-id");
                var start = Math.min(id, that.prevCheckboxID);
                var end = Math.max(id, that.prevCheckboxID);

                if(this_id >= start && this_id <= end){
                    $(this).find('.entry-checkbox').prop("checked", that.prevChecked);
                }
            });
        }

        this.prevChecked = checked;
        this.prevCheckboxID = id;

        this.#CheckboxChange();
    }

    async #InitializeMenu() {
        var that = this;
        var menu_checkbox = $("#menu-checkbox").on("click", (e) => {

            var all_checked = true;

            $(".entry").each(function(){
                if(!$(this).find('.entry-checkbox').prop("checked")){
                    all_checked = false;
                }
            });

            $(".entry").each(function(){
                $(this).find('.entry-checkbox').prop("checked", !all_checked);
            });

            that.#CheckboxChange();
        });
    }

    async #CheckboxChange(){
        var checked_count = 0;
        var total = $(".entry").length;

        $(".entry").each(function(d){
            var checked = $(this).find('.entry-checkbox').prop("checked");
            if(checked){
                checked_count += 1;
            }
        });

        var menu_checkbox = $("#menu-checkbox");
        if(checked_count == 0){
            menu_checkbox.prop("checked", false);
            menu_checkbox.prop("indeterminate", false);
        }
        else if(checked_count == total){
            menu_checkbox.prop("checked", true);
            menu_checkbox.prop("indeterminate", false);
        }
        else{
            menu_checkbox.prop("indeterminate", true);
        }


        var header = $(".header");
        if(checked_count == 0){
            header.find(".logo").removeClass("hidden");
            header.find(".warning").addClass("hidden");
        } else {
            var warning = String(checked_count) + " entries selected"
            header.find(".logo").addClass("hidden");
            header.find(".warning").removeClass("hidden");
            header.find(".warning").text(warning);
        }


        var trash_button = $("#trash-button");
        var download_button = $("#download-button");
        if(checked_count == 0){
            trash_button.addClass("disabled");
            download_button.addClass("disabled");
        } else {
            trash_button.removeClass("disabled");
            download_button.removeClass("disabled");

        }
    }

    #FromString(content) {
        var isJson = true;
        var object = [];
        try {
            object = JSON.parse(content);
        } catch (e) {
            object = [];
            isJson = false;
        }

        if(isJson){
            if(!Array.isArray(object)){
                object = [];
            }
        }
        else {
            var entries = content.split(/[\r\n;]+/);
            for(var i = 0; i < entries.length; i++){
                var segments = entries[i].split(/[:：]/);;
                var w = segments[0] || "";
                w = w.replace(/\s/g, "");
                var k = segments[1] || "";
                k = k.replace(/\s/g, "");
                var l = segments[2] || "";
                l = l.replace(/\s/g, "");
                object.push({w: w, k:k, l:l});
            }
        }
        return object;
    }
    async #InitializeInput() {
        $("#manual-input").on("keydown", async (e) => {
            if (event.keyCode === 13) {
                e.preventDefault();
                var content = $("#manual-input").val();
                var object = this.#FromString(content);

                var newMemory = this.memory.concat(object);
                await this.storage.SaveMemory(newMemory);
                await this.#ReloadMemory();

                $(".entries").find(".entry").remove();
                this.#PopulateEntries();
                this.#CheckboxChange();
            }
        });
        $("#manual-input").on("paste",async (clipevent) => {
            clipevent.preventDefault();
            var content = (clipevent.originalEvent || clipevent).clipboardData.getData('text/plain');
            var object = this.#FromString(content);

            var newMemory = this.memory.concat(object);
            await this.storage.SaveMemory(newMemory);
            await this.#ReloadMemory();

            $(".entries").find(".entry").remove();
            this.#PopulateEntries();
            this.#CheckboxChange();
        });
    }

    // Initialize TrashButtom
    async #InitializeTrash() {
        $("#trash-button").on('mousedown', function(){
            $(this).addClass('active');
        });
        $("#trash-button").on('mouseup', function(){
            $(this).removeClass('active');
        });
        $("#trash-button").on('animationend', () => {
            $("#trash-button").removeClass('active');
            this.#RemoveSelected();
        });
    }

    async #InitializeDownload(){
        $("#download-txt").on("click", () => this.#DownloadTXT()); 
        $("#download-json").on("click", () => this.#DownloadJSON()); 
        $("#download-anki").on("click", () => this.#DownloadAnki()); 
    }

    // Remove all selected entries from browser memory. Reload the page.
    async #RemoveSelected(){
        // Get list of selected entries
        var selection = await this.#GetSelection();
        var newMemory = [];

        for(var i = 0; i < this.memory.length; i++){
            if(selection.indexOf(String(i)) == -1){
                newMemory.push(this.memory[i]);
            }
        }

        await this.storage.SaveMemory(newMemory);
        await this.#ReloadMemory();

        $(".entries").find(".entry").remove();
        this.#PopulateEntries();
        this.#CheckboxChange();
    }

    async #GetSelection(){
        var selection = [];
        $(".entry").each(function(){
            var id = $(this).attr("data-id");
            if($(this).find('.entry-checkbox').prop("checked")){
                selection.push(id);
            }
        });

        return selection;
    }

    async #GetSelected() {
        var selection = await this.#GetSelection();

        var selected = [];

        for(var i = 0; i < this.memory.length; i++){
            if(selection.indexOf(String(i)) != -1){
                selected.push(this.memory[i]);
            }
        }

        return selected;
    }

    #Download(name, content)
    {
        var link = document.createElement('a');
        link.download = name;
        var blob = new Blob([content], {type: 'text/plain'});
        link.href = window.URL.createObjectURL(blob);
        link.click();
    }
    async #DownloadTXT(){
        var selected = await this.#GetSelected();
        var output = "";
        for(var i = 0; i < selected.length; i++){
            var entry = selected[i];
            output = output + this.storage.EntryToString(entry) + "\r\n";
        }
        this.#Download("ame_memory.txt", output);
    }
    async #DownloadJSON(){
        var selected = await this.#GetSelected();
        var output = JSON.stringify(selected);
        this.#Download("ame_memory.json", output);
    }
    async #DownloadAnki(){
        var selected = await this.#GetSelected();
        var URI = "https://amekanji.com/api/post";

        var xmlHttpRequest = new XMLHttpRequest();
        xmlHttpRequest.open('POST', URI, false);
        xmlHttpRequest.setRequestHeader('Content-Type', 'application/json')

        var request = [];
        for(var i = 0; i < selected.length; i++){
            var kanjiword = selected[i].w || "";
            var kanaword = selected[i].k || "";
            var literal = selected[i].l || "";

            var card = {kanjiword : kanjiword, kanaword : kanaword, literal : literal};
            request.push(card);
        }
        var requestbody = JSON.stringify(request);
        xmlHttpRequest.send(requestbody);

        if(xmlHttpRequest.status != 200){
            var header = $(".header");
            var warning = "Could not connect to AmeKanji :("
            header.find(".logo").addClass("hidden");
            header.find(".warning").removeClass("hidden");
            header.find(".warning").text(warning);
            return;
        }

        var rawresponse = xmlHttpRequest.responseText;
        var response = JSON.parse(rawresponse);

        var uuid = response.UUID;

        var redirectionuri = "https://amekanji.com/" + "get.html?id=" + uuid;
        window.open(redirectionuri);
    }

}


$(async function(){
    popup = new Popup();
    popup.Main();
}); 

