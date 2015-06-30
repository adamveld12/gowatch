# :recycle: Gowatch

Compile and run your app during development time. Gowatch will automatically build and restart your app on compile and app exit,
giving you hot reload like behavior during your dev session.

I wrote this to make things easier on me when I'm iterating and learning to golang, so its pretty bare bones and simplistic.

The file watching is done using [fsnotify](https://github.com/go-fsnotify/fsnotify).

## Usage

Simply go to your desired directory (for example, $GOPATH/src/myBadassGoProject) and run the following:

```sh
  gowatch [options]
```

If you want to run Gowatch from a different directory for whatever reason, you can use the -dir="path/to/my/goproject" option.


### Options:

Note: the option that is set on each argument below is the current default if not passed.

`-dir="."`  A directory path to a go package to build/run/watch.

`-ignore=[]`  A comma delimited list of globs for the file watcher to ignore, right now its more like a file extension glob since that's all I really use it for (ie \*.html or \*.css)

`-onerror=true` Restart automatically if there is an error in either the start up process or the app itself. The app will always restart if a file change is made regardless of this setting.

`-wait=1s` How long to wait before starting up the build after an exit.

`-h|-help` Display these usage instructions.

`-debug=false` Shows debug output, for development use

## License:

[WTFYWPL](https://en.wikipedia.org/wiki/WTFPL)
