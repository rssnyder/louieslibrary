.POSTIX:
.SUFFIXES:

.PHONY: build
build:
				go build -o library -v ./cmd/web
				zip -r library.zip ui library
				rsync -varP library.zip ec2: