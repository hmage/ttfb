# ttfb

Measures TTFB (Time To First Byte), inspired by [ttfb.sh](https://github.com/jaygooby/ttfb.sh), written by @jaygooby and @sandeepraju

## Installation

```
go get github.com/hmage/ttfb
```

## Usage

```
ttfb [-n int] <url> [url ...]

  -n int    Number of times to fetch the url (default 1)
```

Output is coded to resemble output of [ttfb.sh](https://github.com/jaygooby/ttfb.sh)