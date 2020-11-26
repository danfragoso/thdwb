git_revision := $(shell git log -1 --pretty="%H")
git_branch := $(shell git log -1 --pretty="%d")
host_info := $(shell uname -s -r -p)
build_date := $(shell date)

ldflags := -X 'main.gitRevision=$(git_revision)' -X 'main.gitBranch=$(git_branch)' -X 'main.buildTime=$(build_date)' -X 'main.hostInfo=$(host_info)'

all: run clean_assets

run: build_assets
	@go run -ldflags "$(ldflags)" *.go --settings="./settings.json"

build: clean build_assets
	@echo "Building THDWB - 🌭"
	@go build -o thdwb -ldflags "$(ldflags) -s -w" *.go
	@chmod 755 thdwb
	@mkdir ../bin; mv thdwb ../bin/

clean:
	@rm -rf ../bin

test:
	@echo -e "Testing Sauce...\n"
	@go test -v sauce/* | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | GREP_COLOR="01;33" egrep --color=always '\s*[a-zA-Z0-9\-_.]+[:][0-9]+[:]|^'
	@echo -e "\n"

	@echo -e "Testing Mayo...\n"
	@go test -v mayo/* | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | GREP_COLOR="01;33" egrep --color=always '\s*[a-zA-Z0-9\-_.]+[:][0-9]+[:]|^'
	@echo -e "\n"

build_assets:
	@go run -tags=bundleAssets assets/bundler.go

clean_assets:
	@rm assets/icons.go assets/pages.go