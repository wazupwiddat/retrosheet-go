DL_BIN := retrosheet-downloader
DL_SRC := cmd/downloader/main.go
LOAD_BIN := retrosheet-dbloader
LOAD_SRC := cmd/dbloader/main.go
MIG_BIN := retrosheet-migration
MIG_SRC := cmd/schema/main.go

all: build_down build_loader build_migration

build_down: $(DL_SRC)
	go build -o bin/$(DL_BIN) $(DL_SRC)

build_loader: $(LOAD_SRC)
	go build -o bin/$(LOAD_BIN) $(LOAD_SRC)

build_migration: $(MIG_SRC)
	go build -o bin/$(MIG_BIN) $(MIG_SRC)