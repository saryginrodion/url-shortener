docker compose -f tests.docker-compose.yaml down -v
docker compose -f tests.docker-compose.yaml up --build -d --wait
go test -v ./tests/integration
