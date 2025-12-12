# passcomb

Password Combination Generator - a console utility for creating password dictionaries by generating all possible combinations from base passwords.

## Features

- Generate password combinations of size 2-4
- Add extra symbols (!@#$%^&*())
- Positional symbol placement (start, end, between parts)
- Split output files by size
- Interactive console interface (lightweight, vim-style navigation)
- Command line support
- Progress bar and generation statistics
- File path auto-completion

## Installation

```bash
git clone https://github.com/iksnevil/passcomb.git
cd passcomb
go build -o passcomb cmd/passcomb/main.go
```

## Usage

### Interactive Console Mode (Default)

```bash
./passcomb
```

Navigation in console mode:
- `j/k` - Move up/down in lists
- `h/l` - Previous/Next step
- `Space` - Toggle selection
- `Enter` - Confirm/Continue
- `q` - Quit

### Command Line Mode

Basic usage:
```bash
./passcomb -i passwords.txt -o combinations.txt -c 2
```

With extra symbols (short aliases):
```bash
./passcomb -i passwords.txt -o combos.txt -c 3 -s '!@#' -p start,end
```

With extra symbols (long names):
```bash
./passcomb --input passwords.txt --output combos.txt --count 3 --symbols '!@#' --positions start,end
```

With file size limit:
```bash
./passcomb -i passwords.txt -o combos.txt -c 4 -m 50
```

## Command Line Options

- `-i, --input string` - Input file with passwords (required in CLI mode)
- `-o, --output string` - Output file for combinations (required in CLI mode)
- `-c, --count int` - Combination size (2-4) [default: 2]
- `-s, --symbols string` - Extra symbols to use (e.g., '!@#$') [default: none]
- `-p, --positions string` - Symbol positions: start,end,between [default: none]
- `-m, --maxsize int` - Max file size in MB [default: 100]
- `-h, --help` - Show help

## Input File Format

Each line in the input file should contain one password:
```
password1
password2
password3
```

## Examples

If the input file contains:
```
a
b
c
```

With `-size 2` parameter, the output file will contain:
```
aa
ab
ac
ba
bb
bc
ca
cb
cc
```

With `-size 3` parameter:
```
aaa
aab
aac
aba
abb
abc
...
ccc
```

With extra symbols:
```bash
./passcomb -input passwords.txt -output combos.txt -size 2 -symbols '!@' -positions start,end
```

Will create additional combinations:
```
!aa
!ab
...
aa!
ab!
...
```

## Project Structure

```
passcomb/
├── cmd/passcomb/          # Main executable
├── pkg/
│   ├── generator/         # Combination generation core
│   ├── interactive/       # Interactive Console Interface
│   └── cli/              # Command line processing
├── internal/
│   ├── config/           # Application configuration
│   └── progress/         # Progress bar
└── go.mod               # Go module
```

## Testing

```bash
go test ./...
```

## Building

```bash
go build -o passcomb cmd/passcomb/main.go
```

## License

MIT License