# :recycle: gowatch

Lint, test, compile and run your app during development time automatically on save, lint error or test failure. 
Gives hot reload like behavior for you during development time.

## Usage
```sh
  gowatch [options]
```

Options:

  `-l|-lint`
  
    Runs `golint` on the source, restarting the process if the linter returns a warning.
    
  `-t|-test`
  
    Runs `go test` on the package, restarting the process if any tests fail.
    
  -m|-main="main.go"
  
    A .go file where main() is defined. This gets passed to go run.

  `-i|-ignore=false`
  
  `-n|--no-restart=error|exit`
  
    Don't automatically restart the watched package if it ends.
    gowatch will wait for a change in the source files.
    If "error", an exit code of 0 will still restart.
    If "exit", no restart regardless of exit code.

  `-h|--help`
  
    Display these usage instructions.

  `-q|--quiet=true`
  
    Suppress DEBUG messages
