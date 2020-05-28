all:
	go build -o bin/rfe

unittest:
	cd myrsa && go test -v -bench X

bench:
	cd myrsa && go test -v -run X -bench Prime

build:
	go build -o bin/rfe