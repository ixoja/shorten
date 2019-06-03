# shorten
This is a simple URL shortener. It is splint into frontend and backend parts which communicate by gRPC.
Frontend serves a static HTML page and provides a simple HTTP API. 

# Build and run
To build and run I use `make` tool. Simply run `make run_backend` and `make run_frontend`.

If you don't want to use `make`, you might build and run manually:
Build: `go build -v -o bin\shorten`
Run backend: `bin\shorten --mode backend --port :9001`
Run frontend: `bin\shorten --mode frontend --html internal\html --port 9000 --web_url localhost --api_url localhost:9001`

NOTICE: order matters. Run backend first, then run frontend.

# TODOs
- add transactions
- dokerize
- make URLs shorter