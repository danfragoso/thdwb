all:
	@go run *.go http://example.com

run:
	@go run *.go $(url)

debug:
	@go run *.go $(url) debug

build:
	@echo -e "Building THDWB - 🌭"
	@go build -o thdwb -ldflags "-s -w" *.go
	@mv thdwb	bin/thdwb
	@chmod 755 bin/thdwb

test:
	@echo -e "Testing Sauce...\n"
	@go test -v sauce/* | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | GREP_COLOR="01;33" egrep --color=always '\s*[a-zA-Z0-9\-_.]+[:][0-9]+[:]|^'
	@echo -e "\n"

	@echo -e "Testing Mayo...\n"
	@go test -v mayo/* | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/'' | GREP_COLOR="01;33" egrep --color=always '\s*[a-zA-Z0-9\-_.]+[:][0-9]+[:]|^'
	@echo -e "\n"