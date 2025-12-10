# passcomb

Password Combination Generator - a console utility for creating password dictionaries by generating all possible combinations from base passwords.

## Features

- Generate password combinations of size 2-4
- Add extra symbols (!@#$%^&*())
- Positional symbol placement (start, end, between parts)
- Split output files by size
- Interactive TUI interface
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

### Interactive Mode (TUI)

```bash
./passcomb -tui
```

### Command Line Mode

Basic usage:
```bash
./passcomb -input passwords.txt -output combinations.txt -size 2
```

With extra symbols:
```bash
./passcomb -input passwords.txt -output combos.txt -size 3 -symbols '!@#' -positions start,end
```

With file size limit:
```bash
./passcomb -input passwords.txt -output combos.txt -size 4 -maxsize 50
```

## Command Line Options

- `-input string` - Input file with passwords (required in CLI mode)
- `-output string` - Output file for combinations (required in CLI mode)
- `-size int` - Combination size (2-4) [default: 2]
- `-symbols string` - Extra symbols to use (e.g., '!@#$') [default: none]
- `-positions string` - Symbol positions: start,end,between [default: none]
- `-maxsize int` - Max file size in MB [default: 100]
- `-tui` - Use interactive interface
- `-help` - Show help

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
│   ├── tui/              # Terminal User Interface
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