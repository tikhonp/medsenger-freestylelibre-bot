templ_serve:
	@templ generate -watch -proxy=http://localhost:9990

templ:
	@templ generate

pkl_conf:
	@pkl-gen-go pkl/config.pkl --base-path github.com/TikhonP/medsenger-freestylelibre-bot

tailwind_serve:
	@tailwindcss -i view/css/input.css -o public/styles.css --watch

tailwind:
	@tailwindcss -i view/css/input.css -o public/styles.css --minify

deploy: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
deploy:
	docker compose -f compose.prod.yaml up --build

docker_dev: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
docker_dev:
	docker compose up --build
