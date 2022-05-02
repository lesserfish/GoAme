#!/bin/sh

if [ "$1" = "audio" ]; then
    shift;
    PORT=8000 CACHE=/tmp TMP=/tmp node Tools/AudioHandler/app.js Tools/AudioHandler/configuration/404.mp3 $@
elif [ "$1" = "api" ]; then
    shift;
    go run Web/Server/API/main.go Web/Server/API/api.go Web/Server/API/endpoints.go Web/Server/API/middleware.go -db Resources/Database/KanjiKana.sqlite $@
elif [ "$1" = "worker" ]; then
    shift;
    go run Web/Server/Background/Ame/main.go Web/Server/Background/Ame/cleaners.go Web/Server/Background/Ame/workers.go -c Resources/configuration.json $@
fi