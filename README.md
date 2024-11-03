# Project Go Solid Boilerplate

Open source Golang Application boilerplate

1. How do you structure a Golang project for production? | Golang Tutorial

[![IMAGE ALT TEXT](https://i.ytimg.com/vi/I_FKeGXIaMs/hqdefault.jpg)](http://www.youtube.com/watch?v=I_FKeGXIaMs "How do you structure a Golang project for production? | Golang Tutorial")

## MakeFile

Run build make command with tests
```bash
make all
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

Live reload the application:
```bash
make watch
```

Clean up binary from the last build:
```bash
make clean
```
