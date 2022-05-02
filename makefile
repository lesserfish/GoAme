default:
	mkdir -p bin/Console
	mkdir -p bin/Server/API
	mkdir -p bin/Server/Workers
	mkdir -p bin/Server/AudioInterface
	mkdir -p bin/Resources
	go build -o bin/Console/console Console/console.go
	go build -o bin/Server/API/API Web/Backend/API/main.go Web/Backend/API/api.go Web/Backend/API/endpoints.go Web/Backend/API/middleware.go
	go build -o bin/Server/Workers/AmeWorker Web/Backend/Background/Ame/main.go Web/Backend/Background/Ame/cleaners.go Web/Backend/Background/Ame/workers.go
	cp -r Resources/* bin/Resources
	cp -r Tools/AudioHandler/app.js Tools/AudioHandler/package.json Tools/AudioHandler/configuration Tools/AudioHandler/.gitignore bin/Server/AudioInterface

clean:
	rm -rf bin/Console
	rm -rf bin/Resources
	rm -rf bin/Server