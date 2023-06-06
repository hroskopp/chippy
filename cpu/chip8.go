package cpu

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	maxROMSpace  = 0xfff - 0x200
	startAddress = 0x200
	ScreenWidth  = 64
	ScreenHeight = 32
)

var sprites = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

var opcodeTable = map[uint8]func(uint16, *Chip8){
	0x0: func(opcode uint16, c *Chip8) {
		switch uint8(opcode & 0x000f) {
		case 0x0: //0x00E0
			c.Screen = [ScreenHeight][ScreenWidth]uint8{}
		case 0xE: //0x00EE
			c.sp--
			c.pc = c.stack[c.sp]
		}
	},
	// 1nnn
	0x1: func(opcode uint16, c *Chip8) {
		address := opcode & 0x0fff
		c.pc = address
	},
	// 2nnn
	0x2: func(opcode uint16, c *Chip8) {
		address := opcode & 0x0fff
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = address
	},
	// 3xkk
	0x3: func(opcode uint16, c *Chip8) {
		regIdx := uint8((opcode & 0x0f00) >> 8)
		vx := c.registers[regIdx]
		otherValue := uint8(opcode & 0x00ff)
		if vx == otherValue {
			c.pc += 2
		}
	},
	// 4xkk
	0x4: func(opcode uint16, c *Chip8) {
		regIdx := uint8((opcode & 0x0f00) >> 8)
		vx := c.registers[regIdx]
		otherValue := uint8(opcode & 0x00ff)
		if vx != otherValue {
			c.pc += 2
		}
	},
	// 5xy0
	0x5: func(opcode uint16, c *Chip8) {
		vx := c.registers[uint8((opcode&0x0f00)>>8)]
		vy := c.registers[uint8((opcode&0x00f0)>>4)]
		if vx == vy {
			c.pc += 2
		}
	},
	// 6xkk
	0x6: func(opcode uint16, c *Chip8) {
		c.registers[uint8((opcode&0x0f00)>>8)] = uint8(opcode & 0x00ff)
	},
	// 7xkk
	0x7: func(opcode uint16, c *Chip8) {
		c.registers[uint8((opcode&0x0f00)>>8)] += uint8(opcode & 0x00ff)
	},
	0x8: func(opcode uint16, c *Chip8) {
		switch uint8(opcode & 0x000f) {
		case 0x0: // 8xy0
			c.registers[uint8((opcode&0x0f00)>>8)] = c.registers[uint8((opcode&0x00f0)>>4)]
		case 0x1: // 8xy1
			c.registers[uint8((opcode&0x0f00)>>8)] |= c.registers[uint8((opcode&0x00f0)>>4)]
		case 0x2: // 8xy2
			c.registers[uint8((opcode&0x0f00)>>8)] &= c.registers[uint8((opcode&0x00f0)>>4)]
		case 0x3: // 8xy3
			c.registers[uint8((opcode&0x0f00)>>8)] ^= c.registers[uint8((opcode&0x00f0)>>4)]
		case 0x4: // 8xy4
			vx, vy := c.registers[uint8((opcode&0x0f00)>>8)], c.registers[uint8((opcode&0x00f0)>>4)]
			sum := uint16(vx) + uint16(vy)
			carry := uint8(sum >> 8)
			c.registers[uint8((opcode&0x0f00)>>8)] = uint8(0x00ff & sum) // only keep lowest 8 bits of sum
			if carry > 0 {
				c.registers[0xf] = 1
			} else {
				c.registers[0xf] = 0
			}

		case 0x5: // 8xy5
			vx, vy := c.registers[uint8((opcode&0x0f00)>>8)], c.registers[uint8((opcode&0x00f0)>>4)]
			if vx > vy {
				c.registers[0xf] = 1
			} else {
				c.registers[0xf] = 0
			}
			c.registers[uint8((opcode&0x0f00)>>8)] = vx - vy
		case 0x6: // 8xy6
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			c.registers[0xf] = vx & 1
			c.registers[uint8((opcode&0x0f00)>>8)] = vx >> 1
		case 0x7: // 8xy7
			vx, vy := c.registers[uint8((opcode&0x0f00)>>8)], c.registers[uint8((opcode&0x00f0)>>4)]
			if vy > vx {
				c.registers[0xf] = 1
			} else {
				c.registers[0xf] = 0
			}
			c.registers[uint8((opcode&0x0f00)>>8)] = vy - vx
		case 0xE: // 8xyE
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			c.registers[0xf] = (vx >> 7) & 1
			c.registers[uint8((opcode&0x0f00)>>8)] = vx << 1
		}
	},
	// 9xy0
	0x9: func(opcode uint16, c *Chip8) {
		vx, vy := c.registers[uint8((opcode&0x0f00)>>8)], c.registers[uint8((opcode&0x00f0)>>4)]
		if vx != vy {
			c.pc += 2
		}
	},
	// Annn
	0xa: func(opcode uint16, c *Chip8) {
		c.idx = opcode & 0x0fff
	},
	// Bnnn
	0xb: func(opcode uint16, c *Chip8) {
		c.pc = (opcode & 0x0fff) + uint16(c.registers[0])
	},
	// Cxkk
	0xc: func(opcode uint16, c *Chip8) {
		c.registers[uint8((opcode&0x0f00)>>8)] = uint8(opcode&0x00ff) & getRandomNum()
	},
	// Dxyn
	0xd: func(opcode uint16, chip *Chip8) {
		vx, vy := chip.registers[uint8((opcode&0x0f00)>>8)], chip.registers[uint8((opcode&0x00f0)>>4)]
		vx, vy = vx%ScreenWidth, vy%ScreenHeight
		chip.registers[0xf] = 0
		nBytes := opcode & 0x000f
		for r := vy; r < (vy + uint8(nBytes)); r++ {
			spriteRow := chip.memory[chip.idx+uint16(r-vy)]
			if r >= ScreenHeight {
				break
			}
			mask, shiftAmount := uint8(0x80), 7
			for c := vx; c < (vx + 8); c++ {
				if c >= ScreenWidth {
					break
				}
				currentPixel := chip.Screen[r][c]
				newPixel := currentPixel ^ ((spriteRow & mask) >> shiftAmount)
				chip.Screen[r][c] = newPixel
				if currentPixel == 1 && newPixel == 0 {
					chip.registers[0xf] = 1
				}
				mask >>= 1
				shiftAmount--
			}
		}
	},
	0xe: func(opcode uint16, c *Chip8) {
		switch uint8(opcode & 0x000f) {
		case 0xe: // Ex9E
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			if c.KeysPressed[vx] {
				c.pc += 2
			}
		case 0x1: // ExA1
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			if !c.KeysPressed[vx] {
				c.pc += 2
			}
		}
	},
	0xf: func(opcode uint16, c *Chip8) {
		switch uint8(opcode & 0x00ff) {
		case 0x07: // Fx07
			c.registers[uint8((opcode&0x0f00)>>8)] = c.delayTimer
		case 0x0a: // Fx0A
			for keyIdx, isPressed := range c.KeysPressed {
				if isPressed {
					c.registers[uint8((opcode&0x0f00)>>8)] = uint8(keyIdx)
					return
				}
			}
			c.pc -= 2
		case 0x15: // Fx15
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			c.delayTimer = vx
		case 0x18: // Fx18
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			c.soundTimer = vx
		case 0x1e: // Fx1eE
			vx := uint16(c.registers[uint8((opcode&0x0f00)>>8)])
			c.idx += vx
		case 0x29: // Fx29
			vx := uint16(c.registers[uint8((opcode&0x0f00)>>8)])
			c.idx = 5 * vx
		case 0x33: // Fx33
			vx := c.registers[uint8((opcode&0x0f00)>>8)]
			for i := 2; i > -1; i-- {
				c.memory[c.idx+uint16(i)] = vx % 10
				vx = uint8(vx / 10)
			}
		case 0x55: // Fx55
			regIdx := uint8((opcode & 0x0f00) >> 8)
			copy(c.memory[c.idx:], c.registers[:regIdx+1])
		case 0x65: // Fx65
			regIdx := uint8((opcode & 0x0f00) >> 8)
			copy(c.registers[:regIdx+1], c.memory[c.idx:])
		}
	},
}

type Chip8 struct {
	memory      [4096]uint8 // 4 KB of memory
	registers   [16]uint8   // 16 8-bit registers (V0 - VF)
	idx         uint16      // 16-bit index register to store memory addresses
	pc          uint16      // (program counter) holds address of next instruction
	stack       [16]uint16  // stores up to 16 memory addresses
	sp          uint8       // stack pointer
	delayTimer  uint8
	soundTimer  uint8
	Screen      [ScreenHeight][ScreenWidth]uint8 // 64x32-pixel display
	KeysPressed [16]bool
}

func getRandomNum() uint8 {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return uint8(r.Intn(256))
}

func (c *Chip8) loadSprites() {
	copy(c.memory[:len(sprites)], sprites[:])
}

func NewCPU() *Chip8 {
	newChip := &Chip8{}
	newChip.loadSprites()
	return newChip
}

func (c *Chip8) LoadROM(f *os.File) error {
	fileInfo, err := os.Stat(f.Name())
	if err != nil {
		return err
	}
	if fileInfo.Size() > maxROMSpace {
		return fmt.Errorf("the program is %d bytes too big", fileInfo.Size()-maxROMSpace)
	}
	buffer := make([]byte, fileInfo.Size())
	if _, err := f.Read(buffer); err != nil {
		return err
	}
	for i, b := range buffer {
		c.memory[startAddress+i] = b
	}
	c.pc = 0x200 // point to first instruction
	return nil

}

// fetch, decode, and execute
func (c *Chip8) Cycle() {
	opcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc += 2
	opcodeTable[uint8((opcode&0xf000)>>12)](opcode, c)
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.soundTimer > 0 {
		c.soundTimer--
	}
}
