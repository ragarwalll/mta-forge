# MTAForge

MTAForge is a CLI tool that helps you create and manage MTA:SA server resources and modules.

## Installation

You can install MTAForge using the following command:

```bash
go install github.com/ragarwalll/mta-forge@latest
```

## Usage

To use MTAForge, you can run the following command:

```bash
mta-forge --help

Usage:
  mf [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  generate    Generate mta and *.mtaext files
  help        Help about any command

Flags:
  -b, --base-dir string     Specify the base directory where the templates are located (default: current directory)
      --expand-source       Expand source information in logs
  -h, --help                help for mf
  -l, --local               Enable local output format
  -o, --output-dir string   Specify the output directory where the generated files will be saved (default: ./output)
  -v, --verbose             Enable verbose output

Use "mf [command] --help" for more information about a command.
```

## License

MTAForge is licensed under the GPL-3.0 license. See the [LICENSE](LICENSE) file for more details.
