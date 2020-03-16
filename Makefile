#
# Prepare openapi files and run projects in containers.
#
YAML=$(shell find * -name \*yaml)
YAMLSRC=$(shell find openapi -name \*yaml.src)
YAMLGEN=$(patsubst %.yaml.src,%.yaml,$(YAMLSRC))
ECHO_GEN=generated/github.com/teamdigitale/api-starter-kit-go/api/api-types.gen.go generated/github.com/teamdigitale/api-starter-kit-go/api/api-server.gen.go
CHI_GEN=generated/github.com/teamdigitale/api-starter-kit-go/api/api-types.gen.go generated/github.com/teamdigitale/api-starter-kit-go/api/api-server-chi.gen.go

yaml: $(YAMLGEN)

.ONESHELL:
%.yaml: %.yaml.src
	tox -e yamllint -- -d relaxed $<
	tox -e yaml 2>/dev/null --  $< $@

yamllint: $(YAML)
	tox -e yamllint -- $<


%-types.gen.go: openapi/simple.yaml
	oapi-codegen  -package api --generate types -o $@ $<
	sed -i -e 's,*string,string,g' -e 's,*int32,int32,g' -e 's,*time.Time,time.Time,g' $@

%-server.gen.go: openapi/simple.yaml
	oapi-codegen  -package api --generate server,spec -o $@ $<

%-server-chi.gen.go: openapi/simple.yaml
	oapi-codegen  -package api --generate chi-server,spec -o $@ $<


prepare-echo:
	mkdir -p  generated/github.com/teamdigitale/api-starter-kit-go
	cp -rp go-echo/* generated/github.com/teamdigitale/api-starter-kit-go

prepare-chi:
	mkdir -p  generated/github.com/teamdigitale/api-starter-kit-go
	cp -rp go-chi/* generated/github.com/teamdigitale/api-starter-kit-go

echo-gen: prepare-echo $(ECHO_GEN)

chi-gen: prepare-chi $(CHI_GEN)

go-build:
	cd generated/github.com/teamdigitale/api-starter-kit-go
	go mod init github.com/teamdigitale/api-starter-kit-go
	go build

build-chi: clean chi-gen go-build

build-echo: clean echo-gen go-build

run-echo: build-echo
	cd generated/github.com/teamdigitale/api-starter-kit-go && go run main.go

run-chi: build-chi
	cd generated/github.com/teamdigitale/api-starter-kit-go && go run main.go
	

test-echo: build-echo
	cd generated/github.com/teamdigitale/api-starter-kit-go && go test

test-chi: build-chi
	cd generated/github.com/teamdigitale/api-starter-kit-go && go test


clean:
	rm -rf generated

reformat:
	find . -name *.go -exec gofmt -w {} \;
