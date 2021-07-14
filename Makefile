run:
	go run . --source '$(shell pwd)/../scr-web/storage/story' --target '$(shell pwd)/a'
build:
	go build .