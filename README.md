# packagetool

A tool for working with various archive formats used in Konami arcade games

# Usage

```
Usage: packagetool [OPTIONS] FILENAME
List of available options:
  -l    List archive contents
  -o string
        Path to the ouput directory (default "./")
```

To dump the contents of an archive, simply run `packagetool FILENAME`.

## Supported Formats

- BAR (both variants)
- MAR
- QAR
