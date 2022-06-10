build:
	cd cmd/teal && go build .

run:
	cd cmd/teal && ./teal

sqlite:
	cd storage && jet -source=sqlite -dsn="./schema.db" -path=./sqlite

cover:
	go test ./storage -coverprofile cover.out
	go tool cover -html=cover.out
