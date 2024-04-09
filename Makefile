templ_serve:
	@templ generate -watch -proxy=http://localhost:8080

templ:
	@templ generate

pkl_conf:
	@pkl-gen-go pkl/config.pkl --base-path github.com/TikhonP/medsenger-freestylelibre-bot

tailwind_serve:
	@tailwindcss -i view/css/input.css -o public/styles.css --watch

tailwind:
	@tailwindcss -i view/css/input.css -o public/styles.css --minify

