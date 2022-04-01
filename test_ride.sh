#!/bin/sh

if [ "$1" = "audio" ]; then
    shift;
    PORT=8000 CACHE=/tmp TMP=/tmp node Tools/AudioHandler/app.js Tools/AudioHandler/configuration/404.mp3 $@
elif [ "$1" = "api" ]; then
    shift;
    go run Web/Backend/API/main.go Web/Backend/API/api.go Web/Backend/API/endpoints.go Web/Backend/API/middleware.go -db Resources/Database/KanjiKana.sqlite $@
elif [ "$1" = "worker" ]; then
    shift;
    go run Web/Backend/Background/Ame/main.go Web/Backend/Background/Ame/cleaners.go Web/Backend/Background/Ame/workers.go -c Resources/configuration.json $@
fi