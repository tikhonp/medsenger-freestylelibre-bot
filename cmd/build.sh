export PATH=$PATH:/usr/local/go/bin
go version
go build -o /home/medsenger/medsenger-freestylelibre-bot/bin/main cmd/server/main.go
go build -o /home/medsenger/medsenger-freestylelibre-bot/bin/fetch_task cmd/fetch_task/main.go
