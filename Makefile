run-example:
	go run example/*.go

test-all:
	go test -v ./

test-cover:
	go test -v ./ --cover

clean-cover:
	rm -rf cover.*

cover-out:
	go test -v -coverprofile cover.out ./

cover-html:
	go tool cover -html cover.out -o cover.html

run-server:
	python3 -m http.server

gen-cover: clean-cover cover-out cover-html run-server

.PHONY: run-example test-all test-cover clean-cover cover-out cover-html gen-cover run-server