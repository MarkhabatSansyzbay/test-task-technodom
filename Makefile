psqlrun:
		docker run --name cont -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -e POSTGRES_DB=db -d postgres
dockersh:
		docker exec -it cont sh
run:
		go run cmd/main.go
