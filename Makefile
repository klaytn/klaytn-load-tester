.PHONY: build

TARGET=./build/bin/klayslave
rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))
SOURCES := $(call rwildcard,.,*.go)

define LD_FLAGS
"-X 'main.Commit=$$(git rev-parse HEAD)'\
 -X 'main.Branch=$$(git rev-parse --abbrev-ref HEAD)'\
 -X 'main.Tag=$$(git name-rev --tags --name-only HEAD)'\
 -X 'main.BuildDate=$$(date)'\
 -X 'main.BuildUser=$$(id -u -n)'"
endef

build: $(TARGET)

$(TARGET): $(SOURCES)
	mkdir -p ./build/bin
	go build -ldflags=${LD_FLAGS} -a -o $(TARGET) ./klayslave
