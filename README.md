# gmxcli: a command line client for GMail -- manage labels, filters using simple YAML config!

[![Go Report Card](https://goreportcard.com/badge/github.com/diwakergupta/gmxcli)](https://goreportcard.com/report/github.com/diwakergupta/gmxcli)


## Installation

Download pre-compiled binaries for Mac, Windows and Linux from the [releases](https://github.com/diwakergupta/gmxcli/releases) page.

## Usage

`gmxcli` is largely self-documenting. Run `gmxcli -h` for a list of supported commands.

```bash
$ gmxcli labels list
...
$ gmxcli filters delete # delete ALL existing filters, use with caution!
...
$ gmxcli filters upload -c gmxcli.yaml # Upload filters defined in the config file
...
```

```yaml
# Sample configuration
filters:
  -
    criteria:
      to: "me"
    action:
      addLabelIds: ["To me"]
      removeLabelIds: ["INBOX", "UNREAD"]
  -
    criteria: ...
 ```
