.SILENT:

PACKAGES=$$(go list ./... | grep -v '/vendor/')
TAG=$$(git describe --tags | cut -c2-)

.PHONY: build
build: jap

.PHONY: test
test:
	go test -cover $(PACKAGES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: preview
preview:
	make -C cmd/jap/ $@

cmd/jap/jap:
	make -C cmd/jap/ jap

jap: cmd/jap/jap *.go
	ln -f cmd/jap/jap jap

.PHONY: clean
clean:
	make -C cmd/jap/ $@
	rm -f jap

deps.svg: *.go
	(   echo "digraph G {"; \
	go list -f '{{range .Imports}}{{printf "\t%q -> %q;\n" $$.ImportPath .}}{{end}}' \
		$$(go list -f '{{join .Deps " "}}' .) .; \
	echo "}"; \
	) | dot -Tsvg -o $@

.PHONY: container
container: jap
	# Get the current Git tag (or latesttag-hash) and drop the first character (so
	# v0.0.5 becomes 0.0.5).
	echo Building container jap:$(TAG)…
	docker build -t jap:$(TAG) .

.PHONY: ecs
ecs:
	make -C ecs/ $<
