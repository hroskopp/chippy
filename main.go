package main

import (
	"flag"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hroskopp/chippy/cpu"
)

type Emulator struct {
	scale int
	chip  *cpu.Chip8
}

func (e *Emulator) Update() error {
	e.processInput()
	e.chip.Cycle()
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	for r := 0; r < cpu.ScreenHeight*e.scale; r++ {
		for c := 0; c < cpu.ScreenWidth*e.scale; c++ {
			if e.chip.Screen[r/e.scale][c/e.scale] > 0 {
				screen.Set(c, r, color.Black)
			} else {
				screen.Set(c, r, color.White)
			}
		}
	}
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return cpu.ScreenWidth * e.scale, cpu.ScreenHeight * e.scale
}

func (e *Emulator) processInput() {
	e.chip.KeysPressed = [16]bool{
		ebiten.IsKeyPressed(ebiten.KeyX),
		ebiten.IsKeyPressed(ebiten.KeyDigit1),
		ebiten.IsKeyPressed(ebiten.KeyDigit2),
		ebiten.IsKeyPressed(ebiten.KeyDigit3),
		ebiten.IsKeyPressed(ebiten.KeyQ),
		ebiten.IsKeyPressed(ebiten.KeyW),
		ebiten.IsKeyPressed(ebiten.KeyE),
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyD),
		ebiten.IsKeyPressed(ebiten.KeyZ),
		ebiten.IsKeyPressed(ebiten.KeyC),
		ebiten.IsKeyPressed(ebiten.KeyDigit4),
		ebiten.IsKeyPressed(ebiten.KeyR),
		ebiten.IsKeyPressed(ebiten.KeyF),
		ebiten.IsKeyPressed(ebiten.KeyV),
	}
}

func main() {
	chippy := &Emulator{chip: cpu.NewCPU()}
	filePath := flag.String("file", "", "Path to ROM file")
	flag.IntVar(&chippy.scale, "scale", 10, "Factor applied to screen dimensions")
	flag.Parse()
	if *filePath == "" {
		log.Fatal("a file must be provided\ncorrect usage: go run main.go -file=path/to/rom")
	}
	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	chippy.chip.LoadROM(f)
	ebiten.SetWindowSize(cpu.ScreenWidth*chippy.scale, cpu.ScreenHeight*chippy.scale)
	ebiten.SetWindowTitle("Chippy: An awesome chip-8 emulator")
	if err := ebiten.RunGame(chippy); err != nil {
		log.Fatal(err)
	}
}
