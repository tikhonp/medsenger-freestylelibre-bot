templ:
	@templ generate -watch -proxy=http://localhost:8080

pkl_conf:
	@pkl-gen-go pkl/config.pkl --base-path github.com/TikhonP/medsenger-freestylelibre-bot

