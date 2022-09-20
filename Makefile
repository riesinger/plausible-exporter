BIN := plausible-exporter

.PHONY: $(BIN)

$(BIN):
	go build -o $(BIN) ./cmd

static:
	CGO_ENABLED=0 go build -o $(BIN) -a -ldflags '-s -w -extldflags "-static"' -tags timetzdata ./cmd

run: $(BIN)
	./$(BIN)
