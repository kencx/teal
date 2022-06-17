build:
	cd cmd/teal && go build .

run:
	cd cmd/teal && ./teal

sqlite:
	cd storage && jet -source=sqlite -dsn="./schema.db" -path=./sqlite

cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out
