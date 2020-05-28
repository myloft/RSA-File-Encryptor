all:
	go build -o rfe

unittest:
	cd myrsa && go test -v -bench X

bench:
	cd myrsa && go test -v -run X -bench Prime