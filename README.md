# 🍇 grape cli

Tiny CLI for running commands on file change. Inspired by [nodemon](https://github.com/remy/nodemon).


## Install 
Download any of the binaries or :

```bash 
go get -u get github.com/noelukwa/grape
```
or 

```bash
go install github.com/noelukwa/grape
```

## Usage

In your working directory, create a `grape.json` file.
  
```go
grape init
```

Create a namespace or run the default with: 

```go
grape run dev
```

Watch and run commands without a config file:

```go
grape on -e ".*go" -c "go run main.go" 
```


## Config Format

```json
{
  "dev": {
    "watch": [
      ".*go"
    ],
    "command": "go run main.go"
  },
  "build": {
    "watch": [
      ".*go"
    ],
    "command": "go build -o main main.go"
  }
}
```

- `watch` is an array of regex patterns to match files to watch
- `command` is the command to run when a file is changed
- `buid,dev` are namespaces to run commands and can be any string




🚧 Under development 🚧
