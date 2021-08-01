test: dockerize
	go run ./cmd/example -h

dockerize:
	go run ./cmd/example dockerize


buildx:
	hx buildx --push

lint:
	husky hook pre-commit
	husky hook commit-msg

