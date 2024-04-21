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

        $(jObject).find(".concept_light-readings").find(".furigana").contents().each(function(index, element){
            if(this.nodeName.toLowerCase() == "span"){
                furigana.push($(this).text());
            }
        });


        var word = $(jObject).find(".concept_light-readings").find(".text").text().replace(/\s+/g, '');

        var kana = "";
        var literal = "";
        
        var counter = 0;
        for(var i = 0; i < furigana.length; i++){
            var f = furigana[i];
            if(f == ""){
                if(word[i]){
                    kana += word[i];
                }
            }
            else {
                if(word[i]){
                    literal += word[i];
                }
                kana += f;
            }
        }

        return({w: word, k: kana, l: literal});
    }

    async #InjectFunctionality(inMemory, eState){
        if(inMemory) {
            var newLink = $("<a href='#' class='concept_light-status_link helper'>Remove from memory</a>");
            newLink.click((e) => this.#DeregisterElement(e, newLink));
            eState.append(newLink);
        }
        else {
            var newLink = $("<a href='#' class='concept_light-status_link helper'>Add to memory</a>");
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
    }

    async #RemoveFromMemory(eState) {
        var element = $(eState).closest(".concept_light.clearfix");
        var reading = this.#GetReadings(element);
        this.memory = this.memory.filter(function(entry){
            return !(entry.w == reading.w && entry.k == reading.k);
        });
        this.storage.SaveMemory(this.memory);
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

