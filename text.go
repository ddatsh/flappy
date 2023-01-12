package main

import (
	_ "embed"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/ttf"
)

//go:embed res/fonts/Flappy.ttf
var embedTtf []byte

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	rectTitle := &sdl.Rect{X: 10, Y: windowHeight / 4, W: windowWidth - 20, H: windowHeight / 2}
	rectPress := &sdl.Rect{X: windowWidth / 4, Y: windowHeight - (windowHeight / 4), W: windowWidth / 2, H: windowHeight / 6}
	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}

	err := drawText(r, "Flappy Gopher", rectTitle, color)
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	err = drawText(r, "press any button to start", rectPress, color)
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()
	return nil
}

func drawText(renderer *sdl.Renderer, text string, rect *sdl.Rect, color sdl.Color) error {
	//path := "res/fonts/Flappy.ttf"
	mem, err := sdl.RWFromMem(embedTtf)
	//font, err := ttf.OpenFont(path, 30)
	font, err := ttf.OpenFontRW(mem, 0, 30)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}

	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return fmt.Errorf("could not render text: %v", err)
	}
	defer surface.Free()

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}

	return renderer.Copy(tex, nil, rect)
}
