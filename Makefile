install:
	sudo chmod -R u+x ./scripts/bash
	@ ./scripts/bash/install.sh
	@ ./scripts/bash/prepare-githooks.sh

lint: 			# lint all the codebase
	golangci-lint run

lint-staged: 	# lint only the staged files
	golangci-lint run --new

debug:
	export GOOGLE_APPLICATION_CREDENTIALS=assets/credentials/firebase.dev.json
	go run cmd/app/main.go

build:
	go build cmd/app/main.go

test:
	go test ./...

gcloud auth:
	export GOOGLE_APPLICATION_CREDENTIALS="${HOME}/go/src/credential.dev.json"
	source ~/.zshrc

run-script:
	go run cmd/script/main.go
