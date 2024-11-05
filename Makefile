BINARY_NAME=procezz
CONFIG_NAME=procezz.conf
PLIST_NAME=procezz.plist

build:
	@echo "Building procezz"
	@go build -ldflags="-s -w" -o tmp/${BINARY_NAME} cmd/procezz/main.go
	@echo "procezz build!"

run: build
	@echo "Starting procezz"
	@./tmp/${BINARY_NAME} &
	@echo "procezz started!"

clean:
	@echo "Cleaning"
	@go clean
	@rm ./tmp/${BINARY_NAME}
	@echo "Cleaned!"

system_install:
	@echo "Adding ${BINARY_NAME} to /usr/local/bin"
	@sudo cp tmp/${BINARY_NAME} /usr/local/bin
	@sudo chmod a+x /usr/local/bin/${BINARY_NAME}
	@echo "Adding ${CONFIG_NAME} to /etc"
	@sudo cp ${CONFIG_NAME} /etc
	@echo "Adding ${PLIST_NAME} to ~/Library/LaunchAgents"
	@cp ${PLIST_NAME} ~/Library/LaunchAgents
	@echo "Loading and Launching ${PLIST_NAME} in ~/Library/LaunchAgents"
	@launchctl load ~/Library/LaunchAgents/${PLIST_NAME}
	@launchctl start ~/Library/LaunchAgents/${PLIST_NAME}

start: run

stop:
	@echo "Stopping procezz"
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped procezz"

restart: stop start

