git pull || exit 1

export PATH=$PATH:/usr/local/go/bin
echo "Building with: $(go version)..."
go build -o /home/medsenger/medsenger-freestylelibre-bot/bin/main cmd/server/main.go || exit 1
go build -o /home/medsenger/medsenger-freestylelibre-bot/bin/fetch_task cmd/fetch_task/main.go || exit 1

sudo supervisorctl restart freestylelibre-medsenger-bot-jobs freestylelibre-medsenger-bot || exit 1

echo "Done! :)"
