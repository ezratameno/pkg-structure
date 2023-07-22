build:
	go build -o pkg-structure ./main.go
	mkdir -p .dist
	mv pkg-structure .dist