SOURCE_COMMIT_SHA := $(shell git rev-parse HEAD)

ENVS := SOURCE_COMMIT=${SOURCE_COMMIT_SHA} COMPOSE_BAKE=true


.PHONY: run dev build-dev prod fprod logs-prod go-to-server-container pkl-gen templ tailwind fetch-task

run: dev

dev:
	${ENVS} docker compose -f compose.yaml up

build-dev:
	${ENVS} docker compose -f compose.yaml up --build

fdev:
	${ENVS} docker compose -f compose.yaml down

prod:
	${ENVS} docker compose -f compose.prod.yaml up --build -d

fprod:
	${ENVS} docker compose -f compose.prod.yaml down

logs-prod:
	${ENVS} docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it --tty agents-freestylelibre-server /bin/bash

pkl-gen:
	docker exec -it --tty agents-freestylelibre-server pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/medsenger-freestylelibre-bot

templ:
	docker exec -it --tty agents-freestylelibre-server templ generate

tailwind:
	tailwindcss -i view/css/input.css -o public/styles.css --minify

fetch-task:
	docker exec -it --tty agents-freestylelibre-server fetch_task
