default: console frontend server 
	
console:
	mkdir -p bin/Console
	go build -o bin/Console/console Console/console.go

frontend:
	cd Web/Frontend && npm install && npm run build
	cp -r Web/Frontend/build/* Resources/public/
server: 
	mkdir -p bin/Server/API
	mkdir -p bin/Server/Workers
	mkdir -p bin/Server/AudioInterface
	mkdir -p bin/Resources
	go build -o bin/Server/API/API Web/Server/API/main.go Web/Server/API/api.go Web/Server/API/endpoints.go Web/Server/API/middleware.go
	go build -o bin/Server/Workers/AmeWorker Web/Server/Background/Ame/main.go Web/Server/Background/Ame/cleaners.go Web/Server/Background/Ame/workers.go
	cp -r Resources/* bin/Resources
	cp -r Tools/AudioHandler/app.js Tools/AudioHandler/package.json Tools/AudioHandler/configuration Tools/AudioHandler/.gitignore bin/Server/AudioInterface


clean:
	rm -rf bin/Console
	rm -rf bin/Resources
	rm -rf bin/Server
	rm -rf Resources/public/*
