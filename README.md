# :recycle: Gowatch

[![Go Report Card](https://goreportcard.com/badge/adamveld12/gowatch)](https://goreportcard.com/report/adamveld12/gowatch)

Compile and run your app during development time. Gowatch will automatically build and restart your app on compile and app exit,
giving you hot reload like behavior during your dev session.

I wrote this to make things easier on me when I'm iterating and learning to golang, so its pretty bare bones and simplistic.

![gowatch demo](/gowatch_demo.gif)

## Usage

Install by doing:
```
go get -u github.com/adamveld12/gowatch
```

Then simply run:
```
  gowatch [options] <path to main package>
```

### Options:

Note: the option that is set on each argument below is the current default if not passed.

`-output=""`the name of the program to output

`-args=""`CLI arguments passed to the app

`-ignore=".git/*,node_modules/*"`  A comma delimited list of globs for the file watcher to ignore, right now its more like a file extension glob since that's all I really use it for (ie \*.html or \*.css)

 `-onexit=true`  If the app should restart on exit, regardless of exit code (beware, buggy)

`-onerror=true` If the app should restart on lint/test/build/non-zero exit code (beware, also buggy)

`-wait=2s` How long to wait before starting up the build after an exit.

 `-test=true` Should tests be run on reload

 `-lint=true` Should the source be linted on reload

`-h|-help` Display these usage instructions.

`-debug=false` Shows debug output, for development use


## Notes

The linter will not stop the app from running if the lint error has a low confidence value (e.g. for missing package level comments). It would be better if a confidence threshold value could be passed.

Process signals (Interrupt, Kill) are not passed to the child process, I ran into a lot of issues with this but I might figure it out in the near future.


## License:

[WTFYWPL](https://en.wikipedia.org/wiki/WTFPL)
