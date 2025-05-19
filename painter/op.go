package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

type Operation interface {
	Do(t screen.Texture, s *State) bool
}

type OperationFunc func(t screen.Texture, s *State) bool

func (f OperationFunc) Do(t screen.Texture, s *State) bool {
	return f(t, s)
}

var UpdateOp = OperationFunc(func(t screen.Texture, s *State) bool {
	drawBackground(t, s)
	drawFigures(t, s)
	return true
})

var WhiteFill = OperationFunc(func(t screen.Texture, s *State) bool {
	s.bgColor = color.White
	return false
})

var GreenFill = OperationFunc(func(t screen.Texture, s *State) bool {
	s.bgColor = color.RGBA{G: 255, A: 255}
	return false
})

func BgRect(rect image.Rectangle) OperationFunc {
	return func(t screen.Texture, s *State) bool {
		s.bgRect = &rect
		return false
	}
}

func Figure(pos image.Point) OperationFunc {
	return func(t screen.Texture, s *State) bool {
		s.figures = append(s.figures, pos)
		return false
	}
}

func Move(delta image.Point) OperationFunc {
	return func(t screen.Texture, s *State) bool {
		s.moveDelta = delta
		return false
	}
}

var ResetOp = OperationFunc(func(t screen.Texture, s *State) bool {
	s.bgColor = color.Black
	s.bgRect = nil
	s.figures = nil
	s.moveDelta = image.Point{}
	t.Fill(t.Bounds(), color.Black, screen.Src)
	return true
})

func drawBackground(t screen.Texture, s *State) {
	t.Fill(t.Bounds(), s.bgColor, screen.Src)
	if s.bgRect != nil {
		t.Fill(*s.bgRect, color.Black, screen.Src)
	}
}

func drawFigures(t screen.Texture, s *State) {
	for _, fig := range s.figures {
		pos := fig.Add(s.moveDelta)
		drawTShape(t, pos)
	}
}

func drawTShape(t screen.Texture, pos image.Point) {
	vertical := image.Rect(pos.X-30, pos.Y-100, pos.X+30, pos.Y+100)
	horizontal := image.Rect(pos.X-100, pos.Y-30, pos.X+30, pos.Y+30)

	t.Fill(vertical, color.RGBA{B: 255, A: 255}, screen.Src)
	t.Fill(horizontal, color.RGBA{B: 255, A: 255}, screen.Src)
}
