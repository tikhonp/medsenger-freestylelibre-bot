export PATH=$PATH:/usr/local/go/bin
go version
go build -o /home/medsenger/ctg-medsenger-bot/bin/main cmd/main.go
go build -o /home/medsenger/ctg-medsenger-bot/bin/fetch_task cmd/fetch_task.go
