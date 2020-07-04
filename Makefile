.POSTIX:
.SUFFIXES:

.PHONY: build
build:
				go build -o library -v ./cmd/web


.PHONY: package
package:
				zip -r library.zip ui library

.PHONY: deploy
deploy:
				rsync -varP library.zip ec2: