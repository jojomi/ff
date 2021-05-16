# ff

Find Folder in a set of given base directories. Useful if you have your projects cluttered around your system.


## Installation

### Download

Pick one of the pre-compiled [released versions](https://github.com/jojomi/ff/releases/). 

Make sure to pick the right binary for your system and put it to any directory in your `$PATH` environment variable.

### Building Yourself

You need a Go environment setup, then:

    go get -u github.com/jojomi/ff

`ff` is built using the current Go compiler version (1.16 at the time of writing), but may work with older versions too.

## Configuration

Config is done in a [`yaml`](http://yaml.org/) file.

The most important key is `paths` where you can specify a list of directories that are used as base directories for searching.

You can find an example config file for your reference at [config.yml.example](config.yml.example).

Default search paths are
* `$HOME/.ff/config.yml` and 
* `/etc/ff/config.yml`.

If you want to use another config, you can specify `--config myfile.yml` when calling `ff`.


## Usage

Get up to date information about how to use the version you downloaded, use `ff --help`.

```
$ ff --help
find folders that match the given fuzzy search pattern

Usage:
  ff [flags]

Flags:
  -c, --config string   config file (default is $HOME/.ff/config.yml)
  -f, --first-only      print only first element
  -h, --help            help for ff
  -v, --verbose         print search details
      --version         print version details
```

Most of the time you will use **`ff proj`** resulting in a list of directories that match your query (fuzzy search is performed).

If you want to change to the first match directly (cd), you can follow up with **`cd "$(!! -f)"`**.