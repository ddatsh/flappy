package main

import (
	"context"
	_ "embed"
	"fmt"

	img "github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

//go:embed res/imgs/background.png
var embedBackground []byte

//go:embed res/imgs/bird_frame_1.png
var embedBird1 []byte

//go:embed res/imgs/bird_frame_2.png
var embedBird2 []byte

//go:embed res/imgs/bird_frame_3.png
var embedBird3 []byte

//go:embed res/imgs/bird_frame_4.png
var embedBird4 []byte

//go:embed res/imgs/pipe.png
var embedPipe []byte

const (
	minPipeDist = 300
)

type scene struct {
	renderer *sdl.Renderer
	bg       *sdl.Texture
	bgt      int32
	bird     *bird
	pipes    *pipes
}

func newScene(r *sdl.Renderer, speed int32, gravity float64) (s *scene, err error) {
	s = &scene{renderer: r}
	mem, err := sdl.RWFromMem(embedBackground)
	s.bg, err = img.LoadTextureRW(r, mem, true)
	//s.bg, err = img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	s.bird = &bird{
		x:       10,
		y:       windowHeight / 2,
		w:       50,
		h:       43,
		gravity: gravity,
	}
	/*	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("res/imgs/bird_frame_%d.png", i)
		frame, errLoad := img.LoadTexture(r, path)
		if errLoad != nil {
			return nil, fmt.Errorf("could not load bird_frame_%d image: %v", i, err)
		}
		s.bird.frames = append(s.bird.frames, frame)
	}*/

	mem, _ = sdl.RWFromMem(embedBird1)
	frame, _ := img.LoadTextureRW(r, mem, true)
	s.bird.frames = append(s.bird.frames, frame)

	mem, _ = sdl.RWFromMem(embedBird2)
	frame, _ = img.LoadTextureRW(r, mem, true)
	s.bird.frames = append(s.bird.frames, frame)

	mem, _ = sdl.RWFromMem(embedBird3)
	frame, _ = img.LoadTextureRW(r, mem, true)
	s.bird.frames = append(s.bird.frames, frame)

	mem, _ = sdl.RWFromMem(embedBird4)
	frame, _ = img.LoadTextureRW(r, mem, true)
	s.bird.frames = append(s.bird.frames, frame)

	s.pipes = &pipes{
		speed: speed,
		pipes: initialPipes(),
	}

	mem, err = sdl.RWFromMem(embedPipe)
	s.pipes.texture, err = img.LoadTextureRW(r, mem, true)
	//s.pipes.texture, err = img.LoadTexture(r, "res/imgs/pipe.png")
	if err != nil {
		return nil, fmt.Errorf("could not load pipe texture: %v", err)

	}

	return s, nil
}

func initialPipes() []*pipe {
	x := int32(windowWidth)
	pp := make([]*pipe, 0)
	for i := 0; i < 4; i++ {
		pp = append(pp, newPipe(x))
		x = x + minPipeDist
	}

	return pp
}

func (s *scene) restart() {
	s.bird.y = windowHeight / 2
	s.bird.dead = false
	s.bird.speed = 0
	s.bird.score = 0
	s.pipes.pipes = initialPipes()
}

func (s *scene) run(ctx context.Context, fps float64) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !s.bird.dead {
				s.update()
				s.draw()
			}

			sdl.Delay(uint32(1000 / fps))
		}
	}
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.bgt = (s.bgt + 1) % 2000

	if s.pipes.hits(s.bird) {
		s.bird.dead = true
	}
}

func (s *scene) draw() error {
	s.renderer.Clear()

	bgRect := &sdl.Rect{X: s.bgt, Y: 0, W: windowWidth, H: windowHeight}
	err := s.renderer.Copy(s.bg, bgRect, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	err = s.bird.draw(s.renderer)
	if err != nil {
		return fmt.Errorf("could not draw bird: %v", err)
	}

	s.pipes.draw(s.renderer)

	if s.bird.dead {
		deadRect := &sdl.Rect{X: 100, Y: windowHeight / 4, W: windowWidth - 200, H: windowHeight / 2}
		drawText(s.renderer, "YOU DIED!", deadRect, sdl.Color{R: 255})
	}

	s.renderer.Present()
	return nil
}
