.PHONY: build

define LD_FLAGS
"-X 'main.Commit=$$(git rev-parse HEAD)'\
 -X 'main.Branch=$$(git rev-parse --abbrev-ref HEAD)'\
 -X 'main.Tag=$$(git name-rev --tags --name-only HEAD)'\
 -X 'main.BuildDate=$$(date)'\
 -X 'main.BuildUser=$$(id -u -n)'"
endef

build:
	mkdir -p ./build/bin
	go build -ldflags=${LD_FLAGS} -a -o ./build/bin/klayslave ./klayslave