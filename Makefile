run:
	go run . --source '$(shell pwd)/../scr-web/storage/story' --target '$(shell pwd)'
build:
	go build .