# go version 1.21.5

# Development

## Build and Run 

1. `go build -o bin/dev-app ./cmd/dev`
2. `nohup ./bin/dev-app &> dev-app.log &`

## Check Running Processes 
`ps aux | grep dev-app` or `pgrep -fl dev-app`

## Killing the Process 
first find PID `pgrep dev-app`
then kill PID `kill <PID>`

# Production

## Build and Run 
1. `go build -o bin/prod-app ./cmd/prod`
2. `nohup ./bin/pro-app &> prod-app.log &`

## Check Running Processes 
`ps aux | grep pro-app` or `pgrep -fl pro-app`

## Killing the Process 
first find PID `pgrep prod-app`
then kill PID `kill <PID>`