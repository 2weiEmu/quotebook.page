main: src/main.go
	@go build -o main src/main.go src/dbStmt.go 

run: src/main.go
	@go run src/main.go src/dbStmt.go

test: src/main.go
	@go build -o main-test src/main.go src/dbStmt.go 
