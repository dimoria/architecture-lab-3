package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

type mockTexture struct {
	size   image.Point
	colors []color.Color
}

func (m *mockTexture) Release()                                                     {}
func (m *mockTexture) Size() image.Point                                            { return m.size }
func (m *mockTexture) Bounds() image.Rectangle                                      { return image.Rectangle{Max: m.size} }
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.colors = append(m.colors, src)
}

type mockScreen struct{}

func (m *mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) { return nil, nil }
func (m *mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{size: size}, nil
}
func (m *mockScreen) NewWindow(*screen.NewWindowOptions) (screen.Window, error) { return nil, nil }

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

func TestLoop(t *testing.T) {
	loop := NewLoop()
	mockScreen := &mockScreen{}
	receiver := &testReceiver{}
	loop.Receiver = receiver

	loop.Start(mockScreen)
	defer loop.StopAndWait()

	loop.Post(WhiteFill)
	loop.Post(UpdateOp)
}
