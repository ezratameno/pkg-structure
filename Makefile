build:
	go build ./app/services/pkg-diff/
	mkdir -p .dist
	mv pkg-diff .dist

run: build
	./.dist/pkg-diff