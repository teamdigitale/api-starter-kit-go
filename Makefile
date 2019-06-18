#
# Prepare openapi files and run projects in containers.
#
YAML=$(shell find * -name \*yaml)
YAMLSRC=$(shell find openapi -name \*yaml.src)
YAMLGEN=$(patsubst %.yaml.src,%.yaml,$(YAMLSRC))
GOGEN=generated/github.com/ioggstream/simple/api/api-types.gen.go generated/github.com/ioggstream/simple/api/api-server.gen.go

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


go-prepare:
	mkdir -p  generated/github.com/ioggstream/simple
	cp -rp go/* generated/github.com/ioggstream/simple

go-build: go-prepare $(GOGEN)
	cd generated/github.com/ioggstream/simple
	go mod init github.com/ioggstream/simple
	go build

run: go-build
	cd generated/github.com/ioggstream/simple && go run main.go

clean:
	rm -rf generated