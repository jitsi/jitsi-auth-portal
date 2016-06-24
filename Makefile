PACKAGES=$$(go list ./... | grep -v '/vendor/')

.PHONEY: build
build: jwtsi

.PHONEY: test
test:
	go test -cover $(PACKAGES)

.PHONEY: vet
vet:
	go vet $(PACKAGES)

cmd/jwtsi/jwtsi:
	make -C cmd/jwtsi/

jwtsi: cmd/jwtsi/jwtsi
	ln -f cmd/jwtsi/jwtsi jwtsi

deps.svg: *.go
	(   echo "digraph G {"; \
	go list -f '{{range .Imports}}{{printf "\t%q -> %q;\n" $$.ImportPath .}}{{end}}' \
		$$(go list -f '{{join .Deps " "}}' .) .; \
	echo "}"; \
	) | dot -Tsvg -o $@
