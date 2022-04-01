#!/bin/sh

if [ "$1" = "audio" ]; then
    PORT=8000 CACHE=/tmp TMP=/tmp node Tools/AudioHandler/app.js Tools/AudioHandler/configuration/404.mp3
elif [ "$1" = "api" ]; then
    go run Web/Backend/API/main.go Web/Backend/API/api.go Web/Backend/API/endpoints.go Web/Backend/API/middleware.go -db Resources/Database/KanjiKana.sqlite
elif [ "$1" = "worker" ]; then
    go run Web/Backend/Background/main.go Web/Backend/Background/cleaners.go Web/Backend/Background/workers.go -c Resources/configuration.json
fi