# chippy

CHIP-8 emulator implemented in Go.

### Installation

1. Clone the repository
2. Install project dependencies

```
go mod download
```

### How To Run

```
go run main.go -file=path/to/rom
```

**Note:** run the following to get information about the program flags

```
go run main.go -h
```

### Keyboard Mapping

| Original Layout | Modern Layout |
| --------------- | ------------- |
| 1               | 1             |
| 2               | 2             |
| 3               | 3             |
| C               | 4             |
| 4               | Q             |
| 5               | W             |
| 6               | E             |
| D               | R             |
| 7               | A             |
| 8               | S             |
| 9               | D             |
| E               | F             |
| A               | Z             |
| 0               | X             |
| B               | C             |
| F               | V             |

### Improvements

- [ ] Add sound
- [ ] Add debugger (to help visualize program execution)
- [ ] Create a CHIP-8 game

### References

- [Cowgod's Chip-8 Technical Reference v1.0](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [BUILDING A CHIP-8 EMULATOR](https://austinmorlan.com/posts/chip8_emulator/)
