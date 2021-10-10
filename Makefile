FILES = $(shell find -name "*.go" -type f 2>/dev/null | sort)
DIRS = $(shell find -name "*.go" -type f 2>/dev/null | xargs -r dirname | sort -u)

main: $(FILES)
	go build -i -o $@ -v

build: fmt main

fmt:
	for x in $(DIRS); do ( cd $$x && go fmt ); done

server: main
	./$< -server

client: main
	./$<

clean:
	-rm -vf main

.PHONY: fmt build
