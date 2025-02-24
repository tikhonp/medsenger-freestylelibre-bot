run: dev

dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
dev:
	docker compose -f compose.yaml up

build-dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
build-dev:
	docker compose -f compose.yaml up --build

prod: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
prod:
	docker compose -f compose.prod.yaml up --build -d

fprod:
	docker compose -f compose.prod.yaml down

logs-prod:
	docker compose -f compose.prod.yaml logs -f -n 100

go-to-server-container:
	docker exec -it agents-freestylelibre /bin/bash

templ-serve:
	templ generate -watch -proxy=http://localhost:9990

templ:
	templ generate

pkl-gen:
	pkl-gen-go pkl/config.pkl --base-path github.com/TikhonP/medsenger-freestylelibre-bot

tailwind-serve:
	tailwindcss -i view/css/input.css -o public/styles.css --watch

tailwind:
	tailwindcss -i view/css/input.css -o public/styles.css --minify
