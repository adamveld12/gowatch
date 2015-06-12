# :recycle: gowatch

Lint, test, compile and run your app during development time automatically on save, lint error or test failure. 
Gives hot reload like behavior for you during development time.

## Usage

```sh
  gowatch [options]
```

Options:

Note: the option that is set on each argument below is the current default if not passed.


 -l|-lint=false

Runs `golint` on the source, restarting the process if the linter returns a warning.


  
-t|-test=false

Runs `go test` on the package, restarting the process if any tests fail.


  
-d|-dir="."

A go package to build.



-i|-ignore=[]

A comma delimited list of globs for the file watcher to  ignore


-n|--norestart=

Don't automatically restart the watched package if it ends.
gowatch will wait for a change in the source files.
If "error", an exit code of 0 will still restart.
If "exit", no restart regardless of exit code.


-h|--help

Display these usage instructions.


-q|--debug=false

Suppress DEBUG messages
