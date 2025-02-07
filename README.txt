FOR START
    go run cmd/main.go --config=./config/config.yaml
OR
    task start
OR
    docker build -t messenger .
    docker run -p 8080:8080 -e CONFIG_PATH=/config/config.yaml messenger
OR
    docker-compose up --build

Then you open it from the project files INDEX.HTML

