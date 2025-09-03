# Project AccountManagementSystem

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the API documentation and Run
```bash
make document_api
make run
```


Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```


error handling notes:
- In case it show "TCP 0.0.0.0:5432 -> 127.0.0.1:0: listen tcp 0.0.0.0:5432: bind: address already in use"
   - first use this cmd on terminal "lsof -i :5432" or 
   - if the above do not show up any result use this cmd on  terminal "sudo lsof -i :5432"
   - then use first try this "kill -9 <pid>" if failed then this "sudo kill -9 <pid>" 
