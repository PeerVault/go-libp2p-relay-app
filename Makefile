build:
	go build -o bin/relay relay.go
	go build -o bin/dialer dialer.go
	go build -o bin/listener listener.go
