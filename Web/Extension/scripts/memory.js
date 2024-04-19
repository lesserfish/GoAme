var MAX_SEG_SIZE = 8000
var Dictionary = []

function ByteSize(str)
{
    return (new TextEncoder().encode(str)).length;
}

function ChunkName(chunk_id) {
    return "chunk_" + chunk_id.toString();
}

function EntryToString(entry)
{
    var output = entry['k'] + ":" + entry['h'];
    return output;
}

function LoadChunk(chunk) {
    var entries = chunk.split('\n');
    for(var i = 0; i < entries.length; i++)
    {
        entry = entries[i];
        if(entry.length == 0){
            continue
        }
        var k = "";
        var h = "";
        var components = entry.split(':');
        k = components[0];
        if(components.length > 1){
            h = components[1];
        }
        var word = {
          k: k,
          h: h,
        }
        Dictionary.push(word)
    }
}
async function SaveChunk(chunk, chunk_id){
    var chunk_name = ChunkName(chunk_id)
    var map = {};
    map[chunk_name] = chunk; 
    await browser.storage.sync.set(map);
}
async function LoadDictionary() {
    var promise = browser.storage.sync.get();
    promise.then( 
        async function(storage) {
            if("chunk_count" in storage){
                for(var chunk_id = 1; chunk_id <= storage['chunk_count']; chunk_id++)
                {
                    var chunk_name = ChunkName(chunk_id);
                    if(!chunk_name in storage)
                    {
                        console.error("CRITICAL ERROR: Attempted to get non-existent chunk");
                    }
                    LoadChunk(storage[chunk_name])
                }
            } else {
                SaveDictionary();
            }
        },
        async function(e) {
            console.error("Error: " + e);
        });
    return await promise
}


async function SaveDictionary() {
    var chunk_count = 1;
    var current_chunk = ""
    for(var i = 0; i < Dictionary.length; i++)
    {
        var entry = EntryToString(Dictionary[i])
        if(ByteSize(current_chunk) + ByteSize("\n") + ByteSize(entry) > MAX_SEG_SIZE)
        {
            SaveChunk(current_chunk, chunk_count);
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
    SaveChunk(current_chunk, chunk_count);
    await browser.storage.sync.set({chunk_count: chunk_count})
}
