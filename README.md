# Chip8

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Chip-8 interpreter written in Go.

> [Chip-8](https://en.wikipedia.org/wiki/CHIP-8) is a simple, interpreted, programming language which was first used on some do-it-yourself computer systems in the late 1970s and early 1980s. With only 4KB of memory and 36 instructions it is the perfect system to build when getting into emulator building.

## Requirements

- [Go 1.18 or higher](https://go.dev/dl/)

## Installation

Clone the repository and install on your machine.

```bash
git clone https://github.com/Smelton01/chip8
cd chip8 && go install cmd/chip8/main.go
./chip8
```

Select your favorite game/rom from the menu and start playing.

## Keyboard Layout:

The CHIP-8 machine uses a hexadecimal keypad with 16 keys, labelled 0 through F, arranged in a 4x4 grid and mapped to a QWERTY layout keyboard as follows.

### Chip8 Keypad:

|     |     |     |     |
| --- | --- | --- | --- |
| 1   | 2   | 3   | C   |
| 4   | 5   | 6   | D   |
| 7   | 8   | 9   | E   |
| A   | 0   | B   | F   |

### Emulator Keyboard Mapping:

|     |     |     |     |
| --- | --- | --- | --- |
| 1   | 2   | 3   | 4   |
| Q   | W   | E   | R   |
| A   | S   | D   | F   |
| Z   | X   | C   | V   |

## Testing

Run the following command to run the unit tests for the packages.

```bash
go test -v ./...
```

## License

This project is open source and available under the [MIT License](LICENSE).
