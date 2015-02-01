# ix

CLI tool to explore closed GitHub issue for a repository by time frame, labels and assignee.

```
$ ix -repo mhfs/ix -t <GH_TOKEN>

#1 - 2015-01-29 - Fixed that weird bug by @mhfs (labe1, label2)
#2 - 2015-01-27 - Added another amazing feature by @mhfs
```

## Installation

Only available from source at the moment. Binary releases planned for future versions.

```
go get github.com/mhfs/ix
cd $GOPATH/src/github.com/mhfs/ix
go install
$GOPATH/bin/ix help
```

## Authentication

ix depends on a oauth token from GitHub. You can provide it via a `--token`/`-t` options or set a
`GH_TOKEN` environment variable.

To generate a new token, go to https://github.com/settings/applications.

## Usage

```
NAME:
   ix - cli to explore closed GitHub issue for a repository by time frame, labels and assignee

USAGE:
   ix [options]

EXAMPLES:
   ix --repo mhfs/ix --since 2015-01-01
   ix --repo mhfs/ix --assignee mhfs
   ix --repo mhfs/ix --label bug

VERSION:
   0.0.1

AUTHOR:
  Marcelo Silveira - <marcelo@mhfs.com.br>

OPTIONS:
   --repo, -r 					GitHub repository to analyze e.g. mhfs/ix
   --since, -s '2015-01-26'			list issues since given date, inclusive
   --label, -l '--label option --label option'	label to process, defaults to all
   --assignee, -a 				filter results by assignee
   --token, -t 					oauth token. defaults to GH_TOKEN env var. [$GH_TOKEN]
   --help, -h					show help
   --version, -v				print the version
```

## License

Released under the MIT License. See the [LICENSE][license] file for further details.

[license]: https://github.com/mhfs/ix/blob/master/LICENSE
