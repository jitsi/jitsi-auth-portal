PACKAGES=$$(go list ./... | grep -v '/vendor/')

.PHONEY: build
build: jap

.PHONEY: test
test:
	go test -cover $(PACKAGES)

.PHONEY: vet
vet:
	go vet $(PACKAGES)

.PHONEY: preview
preview:
	make -C cmd/jap/ $@

cmd/jap/jap:
	make -C cmd/jap/ jap

jap: cmd/jap/jap *.go
	ln -f cmd/jap/jap jap

deps.svg: *.go
	(   echo "digraph G {"; \
	go list -f '{{range .Imports}}{{printf "\t%q -> %q;\n" $$.ImportPath .}}{{end}}' \
		$$(go list -f '{{join .Deps " "}}' .) .; \
	echo "}"; \
	) | dot -Tsvg -o $@
