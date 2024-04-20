class AmeExtension {
    async Main() {
        var storage = new Storage();
        var memory = await storage.LoadMemory();
        this.storage = storage;
        this.memory = memory;

        this.#Go();
        this.#CreateObserver();
        
    }

    async #Go() {
        $(".concept_light.clearfix").each((index, element) => this.#HandleEntry(index, element));
    }

    async #HandleEntry(index, element)
    {
        if($(element).parent().attr("class") == "names") {
            return;
        }
        var eState = $(element).find(".concept_light-status");
        var reading = this.#GetReadings(element);

        var inMemory = this.memory.some(function(entry) {
            return (entry.w == reading.w && entry.k == reading.k);
        });

        this.#InjectFunctionality(inMemory, eState);
    }

    #GetReadings(jObject)
    {
        var furigana = [];

        $(jObject).find(".concept_light-readings").find(".furigana").find(".kanji").each(function(index, element){
            furigana.push($(this).text());
        });

        var aWord = [];
        var aKana = [];

        $(jObject).find(".concept_light-readings").find(".text").contents().each(function(index, element){
            var content = $(this).text().replace(/\s/g, "");
            if(this.nodeName.toLowerCase() == "span"){
                aWord.push(content);
                aKana.push(content);
            }
            else {
                aWord.push(content);
                aKana.push(furigana.shift() || "");
            }
        });

        var word = aWord.join("");
        var kana = aKana.join("");

        return({w: word, k: kana});
    }

    async #InjectFunctionality(inMemory, eState){
        if(inMemory) {
            var newLink = $("<a href='#' class=''>Remove from memory</a>");
            newLink.click((e) => this.#DeregisterElement(e, newLink));
            eState.append(newLink);
        }
        else {
            var newLink = $("<a href='#' class=''>Add to memory</a>");
            newLink.click((e) => this.#RegisterElement(e, newLink));
            eState.append(newLink);
        }

    }
    async #RegisterElement(e, linkObject){
        var eParent = $(linkObject).parent();
        this.#AddToMemory(eParent);
        $(linkObject).remove();
        this.#InjectFunctionality(true, eParent);
        e.preventDefault();
    }
    async #DeregisterElement(e, linkObject){
        var eParent = $(linkObject).parent();
        this.#RemoveFromMemory(eParent);
        $(linkObject).remove();
        this.#InjectFunctionality(false, eParent);
        e.preventDefault();    
    }

    async #AddToMemory(eState) {
        var entry = $(eState).closest(".concept_light.clearfix");
        var reading = this.#GetReadings(entry);
        this.memory.push(reading);
        this.storage.SaveMemory(this.memory);
        console.log(this.memory);
    }

    async #RemoveFromMemory(eState) {
        var element = $(eState).closest(".concept_light.clearfix");
        var reading = this.#GetReadings(element);
        this.memory = this.memory.filter(function(entry){
            return !(entry.w == reading.w && entry.k == reading.k);
        });
        this.storage.SaveMemory(this.memory);
        console.log(this.memory);
    }

    async #CreateObserver(){

        // Select the node that you want to observe for mutations
        const targetNode = document.body; // You can choose any parent node where "primary" div might be added
        const config = { childList: true, subtree: true };
        const observer = new MutationObserver((m, o) => this.#HandleMutation(m, o));

        // Start observing the target node for configured mutations
        observer.observe(targetNode, config);
    }

    async #HandleMutation(mutationsList, observer) {
        // Loop through the list of mutations
        for(const mutation of mutationsList) {
            // Check if nodes were added
            if (mutation.type === 'childList') {
                // Loop through added nodes
                mutation.addedNodes.forEach(node => {
                    // Check if the added node is a div with the id "main_results"
                    if (node.nodeType === 1 && node.id === 'main_results') {
                        this.#Go();
                    }
                });
            }
        }
    }
}

$(async function(){
    ame = new AmeExtension();
    ame.Main();
}); 

