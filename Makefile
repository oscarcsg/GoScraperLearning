Name = Go Scraper Learning
Author = oscarcsg
Description = This project is used to learn how to web scrap in Go using goquery and databases.

Version = dev-0.5.2

# BuildTime = $(shell date +"%Y-%m-%d %H:%M:%S")   ---   -X 'main.BuildTime=$(BuildTime)'
# CommitID = $(shell git rev-parse --short HEAD)   ---   -X 'main.CommitID=$(CommitID)'


LDFLAGS = -ldflags="-X 'main.Name=$(Name)' -X 'main.Author=$(Author)' -X 'main.Description=$(Description)' -X 'main.Version=$(Version)'"

BUILD_OUT = build/scraper-$(Version)

MAIN_FUNC = cmd/scraper/main.go

.PHONY: build run

.SILENT:
build:
	go build $(LDFLAGS) -o $(BUILD_OUT) $(MAIN_FUNC)

run:
	go run $(LDFLAGS) $(MAIN_FUNC)

run-binary:
	./$(BUILD_OUT)