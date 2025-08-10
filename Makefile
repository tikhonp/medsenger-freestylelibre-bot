SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container templ tailwind fetch-task build-prod-image update-deps

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

fdev:
	${ENVS} docker compose -f compose.yaml down

prod:
	docker compose -f compose.test-prod.yaml up --build -d

fprod:
	docker compose -f compose.test-prod.yaml down

logs-prod:
	docker compose -f compose.test-prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty agents-freestylelibre-server /bin/sh

templ:
	docker exec -it --tty agents-freestylelibre-server templ generate

tailwind:
	tailwindcss -i view/css/input.css -o public/styles.css --minify

fetch-task:
	docker exec -it --tty agents-freestylelibre-server fetch_task

build-prod-image:
	docker buildx build --build-arg SOURCE_COMMIT="${SOURCE_COMMIT_SHA}" --target server-prod -t docker.telepat.online/agents-freestylelibre-image:server-latest .
	docker buildx build --build-arg SOURCE_COMMIT="${SOURCE_COMMIT_SHA}" --target worker-prod -t docker.telepat.online/agents-freestylelibre-image:worker-latest .

update-deps:
	docker exec -it --tty agents-freestylelibre-server go get -u ./...
