build:
	go build -o=./tmp/gna cmd/gna/main.go

run: build
	./tmp/gna -debug

logs:
	tail -f ${HOME}/.cache/gimmenews/logs/client.log
