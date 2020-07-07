.POSTIX:
.SUFFIXES:

.PHONY: build
build:
				go build -o dist/library -v ./cmd/web


.PHONY: package
package:
				zip -r dist/library.zip ui library

.PHONY: deploy
deploy:
				rsync -varP dist/library.zip linode:
