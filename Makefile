.PHONY: setup teardown build clean

include .env

teardown: clean
	docker-compose down --remove-orphans
	docker volume prune -f

setup: build
	docker-compose up -d

build: clean
	GOOS=linux GOARCH=amd64 go build
	mkdir -pv ${BIN_DIR}
	mv ${EXPORTER_BIN} ${BIN_DIR}

clean:
	rm -rf ${BIN_DIR}
