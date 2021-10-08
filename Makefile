PATH_TO_FILE=./data-to-load/matchups.json
ifneq ("$(wildcard $(PATH_TO_FILE))","")
else
$(error 'File ./data-to-load/matchups.json is missing.Visit README.md')
endif

node-build:
	(cd node-solution && yarn build)

node-run:
	(cd node-solution && node dist/index.js)

go-build:
	(cd go-solution && CGO_ENABLED=0	go build -ldflags="-w -s" -o dist/etl main.go)

go-run:
	(cd go-solution && dist/etl)

go-lint:
	(cd go-solution && golangci-lint run)

docker-start:
	docker-compose up -d

docker-stop:
	docker-compose down

remove-files: 
	@echo 'Deleting built artifacts/postgres data/error logs ${FILE_EXISTS}'
	rm -rf errors
	rm -rf dist
	rm -rf node-solution/dist
	rm -rf go-solution/dist

clean: docker-stop remove-files
