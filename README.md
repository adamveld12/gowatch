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


### Options:

Note: the option that is set on each argument below is the current default if not passed.

`-args=""`CLI arguments passed to the app

`-ignore=[]`  A comma delimited list of globs for the file watcher to ignore, right now its more like a file extension glob since that's all I really use it for (ie \*.html or \*.css)

 `-onexit=true`  If the app should restart on exit, regardless of exit code

`-onerror=true` If the app should restart on lint/test/build/non-zero exit code

`-wait=1s` How long to wait before starting up the build after an exit.

 `-test=false` Should tests be run on reload

 `-lint=true` Should the source be linted on reload

`-h|-help` Display these usage instructions.

`-debug=false` Shows debug output, for development use

## License:

[WTFYWPL](https://en.wikipedia.org/wiki/WTFPL)
