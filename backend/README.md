# Pre
go mod init test
go get github.com/gin-gonic/gin
go get github.com/lib/pq

# local testing
running program 
go run main.go

# to be used in docker 
to build binary file (prod)
go build main.go

This will create a binary "main" which should be ran on container start
