Name = Go Scraper Learning
Author = oscarcsg
Description = Proyecto para aprender a hacer web scrapping en Go usando goquery.

Version = dev-0.1.0

BuildTime = $(shell date +"%Y-%m-%d %H:%M:%S")
CommitID = $(shell git rev-parse --short HEAD)


LDFLAGS = -ldflags="-X 'main.Name=$(Name)' -X 'main.Author=$(Author)' -X 'main.Description=$(Description)' -X 'main.Version=$(Version)' -X 'main.BuildTime=$(BuildTime)' -X 'main.CommitID=$(CommitID)'"

BUILD_OUT = build/scraper

MAIN_FUNC = cmd/scraper/main.go

.PHONY: build run

.SILENT:
build:
	go build $(LDFLAGS) -o $(BUILD_OUT) $(MAIN_FUNC)

run:
	go run $(LDFLAGS) $(MAIN_FUNC)

run-binary:
	./$(BUILD_OUT)