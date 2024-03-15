const express = require('express');
const axios = require('axios');
const dotenv = require('dotenv');
const fs = require('fs');
const http = require('http')
const crypto = require('crypto');
const { time } = require('console');


dotenv.config();
const AudioHandler = express();
AudioHandler.use(express.json())

const port = process.env.PORT
const cache = process.env.CACHE
const tmpdir = process.env.TMP
const FOFFILE = process.argv.slice(2)[0]
//"./configuration/404.mp3"

var ERRMD5 = ""

// Generate 404.mp3 hash for future comparisons

async function SetUp(){
    ERRMD5 = await checksumFile(FOFFILE)
}
SetUp()


// Beggining of functions
AudioHandler.get('/', async(req, res) => {
    var kana = req.query.kana
    var kanji = req.query.kanji
    
    if(typeof(kana) !== 'string' || !kana instanceof String)
    {
        res.sendStatus(400)
        return
    } 
    
    if(typeof(kanji) !== 'string' || !kanji instanceof String){
        kanji = ""
    }
    
    console.log("[server]: Received request for " + kana + " : " + kanji + ".")

    var filepath = await GetFilepath(kana, kanji)

    if(filepath == ""){
        res.sendStatus(404)
        console.log("[server]: Failed to find file: " + filepath)
        return
    }

    res.sendFile(filepath, (err) => {
        if(err) {
            res.sendStatus(404)
            console.log("[server]: Failed to send file " + filepath + ". " + err.message + ".")
        } else {
            console.log("[server]: Sent file: " + filepath + ".")
        }
    })
})
AudioHandler.listen(port, (err) => {
    if(err){
        console.log(err)
    }
  console.log(`[server]: Server is running at http://localhost:${port}`);
});


async function GetFilepath(kana, kanji)
{
    var fullpath = cache + "/" + GetFilename(kana)
    
    if(fs.existsSync(fullpath)){
        console.log("[server]: File found locally!")
        return fullpath
    }
    
    console.log("[server]: File not found locally! Attempting third party download.")
    
    var err = await TryDownloadFile(kana, kanji, fullpath)

    if(err){
        console.log("[server]: " + err)
        return ""
    }
    return fullpath
}

function GetFilename(kana)
{
    var out = kana + ".mp3"
    return out
}

async function TryDownloadFile(kana, kanji, fullpath)
{
    tmppath = tmpdir + "/" + "ah_tmp_" + String(Math.floor(Math.random() * 10000000)) + ".mp3"
    var URI = ""
    if(kanji == ""){
        URI = `http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kana=${kana}`
    } else {
       URI = `http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kana=${kana}&kanji=${kanji}`
    }


    var total_size = await downloadFile(URI, tmppath)
    if(total_size <= 0)
    {
        return new Error("[server]: Empty file!")
    }

    isValid = await CheckValidity(tmppath)

    if(isValid){
        fs.copyFileSync(tmppath, fullpath)
        return null
    }
    else {
        return new Error("[server]: File "+ tmppath+" is not valid!")
    }
    
}

async function CheckValidity(file)
{
    filehash = await checksumFile(file)
    if(filehash === ERRMD5){
        return false
    }
    return true
}

async function checksumFile(path) {
    return new Promise((resolve, reject) => {
    fs.createReadStream(path).
    pipe(crypto.createHash('sha1').setEncoding('hex')).
    on('finish', function () {
        resolve(this.read())
    })
  }).then(function(output) {
      return output
  })
}


async function downloadFile(URI, tmppath) {
    console.log("Downloading " + URI);

    const writer = fs.createWriteStream(tmppath);

    const response = await axios({
        method: 'get',
        url: URI,
        responseType: 'stream',
        headers: {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
        }
    });

    response.data.pipe(writer);

    return new Promise((resolve, reject) => {
        writer.on('finish', () => {
            console.log('File downloaded successfully to ' + tmppath);
            resolve();
        });

        writer.on('error', (err) => {
            reject(err);
        });
    });
}
