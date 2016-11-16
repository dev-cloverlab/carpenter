# carpenter

Carpenter is a tool to manage DB schema and data

## Description

Carpenter has three sub commands.

- export
    - `export` command is export table structure as JSON string
- build
    - `build` command is migrate table from JSON files
- seed
    - `seed` command is seed table data from CSV files

## Usage

```bash
NAME:
   carpenter - Carpenter is a tool to manage DB schema and data

USAGE:
   carpenter [global options] command [command options] [arguments...]
   
VERSION:
   0.1.0
   
AUTHOR(S):
   hatajoe <hatanaka@cloverlab.jp> 
   
COMMANDS:
     export   Export table structure as JSON string
     build    Build(Migrate) table from specified JSON string
     seed     Seed CSV to table
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose, --vv                show verbose output (default off)
   --dry-run                      execute as dry-run mode (default off)
   --schema value, -s value       database name (required)
   --data-source value, -d value  data source name like '[username[:password]@][tcp[(address:port)]]' (required)
   --help, -h                     show help
   --version, -v                  print the version
```

### export

```bash
NAME:
   carpenter export - Export table structure as JSON string

USAGE:
   carpenter export [command options] [arguments...]

OPTIONS:
   --pretty, -p  show pretty output (default off)
```

### build

```bash
NAME:
   carpenter build - Build(Migrate) table from specified JSON string

USAGE:
   carpenter build [command options] [arguments...]

OPTIONS:
   --dir value, -d value  path to JSON file directory (required)
```

### seed

```bash
NAME:
   carpenter seed - Seed CSV to table

USAGE:
   carpenter seed [command options] [arguments...]

OPTIONS:
   --dir value, -d value  path to CSV file directory (required)
```

NOTICE:

- All tables require an id column
- If you include a line break, please enclose it in double quotation marks
- Please do not put double quotes in double quotes

## Install

To install, use `go get`:

```bash
$ go get -d github.com/dev-cloverlab/carpenter/cmd/carpenter
```

## Contribution

1. Fork ([https://github.com/dev-cloverlab/carpenter/cmd/carpenter/fork](https://github.com/dev-cloverlab/carpenter/cmd/carpenter/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[@hatajoe](https://twitter.com/hatajoe)

## Licence

MIT
