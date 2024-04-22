var MAX_SEG_SIZE = 7500

class Storage {
    #ByteSize(str)
    {
        return (new TextEncoder().encode(str)).length;
    }

    #ChunkName(chunk_id) {
        return "chunk_" + chunk_id.toString();
    }

    EntryToString(entry)
    {
        var output = entry['w'] + ":" + entry['k'] + ":" + entry['l'];
        return output;
    }

    #LoadChunk(chunk, memory) {
        var entries = chunk.split('\n');
        for(var i = 0; i < entries.length; i++)
        {
            var entry = entries[i];
            if(entry.length == 0){
                continue
            }
            var w = "";
            var k = "";
            var l = "";
            var components = entry.split(':');
            w = components[0];
            if(components.length > 1){
                k = components[1];
            }
            if(components.length > 2){
                l = components[2];
            }
            var word = {
              w: w,
              k: k,
              l: l
            }
            memory.push(word)
        }
    }
    async #SaveChunk(chunk, chunk_id){
        var chunk_name = this.#ChunkName(chunk_id)
        var map = {};
        map[chunk_name] = chunk; 
        await browser.storage.sync.set(map);
    }

    async LoadMemory() {
        var memory = [];
        try {
            var storage = await browser.storage.sync.get();
            if ("chunk_count" in storage) {
                for (var chunk_id = 1; chunk_id <= storage['chunk_count']; chunk_id++) {
                    var chunk_name = this.#ChunkName(chunk_id);
                    if (!(chunk_name in storage)) {
                        console.error("CRITICAL ERROR: Attempted to get non-existent chunk");
                    }
                    await this.#LoadChunk(storage[chunk_name], memory);
                }
            } else {
                this.SaveMemory(memory);
            }
            return memory; // Return the populated memory array
        } catch (e) {
            console.error("Error: " + e);
            return []; // Return an empty array in case of error
        }
    }

    async SaveMemory(memory) {
        var chunk_count = 1;
        var current_chunk = ""
        for(var i = 0; i < memory.length; i++)
        {
            var entry = this.EntryToString(memory[i])
            if(this.#ByteSize(current_chunk) + this.#ByteSize("\n") + this.#ByteSize(entry) > MAX_SEG_SIZE)
            {
                this.#SaveChunk(current_chunk, chunk_count);
                current_chunk = entry;
                chunk_count += 1;
            } else {
                if(current_chunk.length == 0){
                    current_chunk = entry;
                } else {
                    current_chunk = current_chunk + "\n" + entry;
                }
            }
        }
        this.#SaveChunk(current_chunk, chunk_count);
        await browser.storage.sync.set({chunk_count: chunk_count})
    }
}
