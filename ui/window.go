package ui

import (
	"image"
	"image/color"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/imageutil"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Visualizer struct {
	Title         string
	Debug         bool
	OnScreenReady func(s screen.Screen)

	w    screen.Window
	tx   chan screen.Texture
	done chan struct{}

	sz  size.Event
	pos image.Rectangle

	//координати центру нашої фігури (в нас був 19 варіант)
	figureX int
	figureY int
}

func (pw *Visualizer) Main() {
	pw.tx = make(chan screen.Texture)
	pw.done = make(chan struct{})
	pw.pos.Max.X = 200
	pw.pos.Max.Y = 200

	//центр за замовчуванням
	pw.figureX = 400
	pw.figureY = 400

	driver.Main(pw.run)
}

func (pw *Visualizer) Update(t screen.Texture) {
	pw.tx <- t
}

func (pw *Visualizer) run(s screen.Screen) {
	w, err := s.NewWindow(&screen.NewWindowOptions{

		//розміри вікна
		Width:  800,
		Height: 800,

		Title: pw.Title,
	})
	if err != nil {
		log.Fatal("Failed to initialize the app window:", err)
	}
	defer func() {
		w.Release()
		close(pw.done)
	}()

	if pw.OnScreenReady != nil {
		pw.OnScreenReady(s)
	}

	pw.w = w

	events := make(chan any)
	go func() {
		for {
			e := w.NextEvent()
			if pw.Debug {
				log.Printf("new event: %v", e)
			}
			if detectTerminate(e) {
				close(events)
				break
			}
			events <- e
		}
	}()

	var t screen.Texture

	for {
		select {
		case e, ok := <-events:
			if !ok {
				return
			}
			pw.handleEvent(e, t)

		case t = <-pw.tx:
			w.Send(paint.Event{})
		}
	}
}

func detectTerminate(e any) bool {
	switch e := e.(type) {
	case lifecycle.Event:
		if e.To == lifecycle.StageDead {
			return true // Window destroy initiated.
		}
	case key.Event:
		if e.Code == key.CodeEscape {
			return true // Esc pressed.
		}
	}
	return false
}

func (pw *Visualizer) handleEvent(e any, t screen.Texture) {
	switch e := e.(type) {

	case size.Event: // Оновлення даних про розмір вікна.
		pw.sz = e

	case error:
		log.Printf("ERROR: %s", e)

	case mouse.Event:
		if t == nil {
			// TODO: Реалізувати реакцію на натискання кнопки миші.
			if e.Button == mouse.ButtonLeft {
				pw.figureX = int(e.X)
				pw.figureY = int(e.Y)
				pw.w.Send(paint.Event{}) //реалізація натисканн на ліву кнопку
			}
		}

	case paint.Event:
		// Малювання контенту вікна.
		if t == nil {
			pw.drawDefaultUI()
		} else {
			// Використання текстури отриманої через виклик Update.
			pw.w.Scale(pw.sz.Bounds(), t, t.Bounds(), draw.Src, nil)
		}
		pw.w.Publish()
	}
}

func (pw *Visualizer) drawDefaultUI() {
	pw.w.Fill(pw.sz.Bounds(), color.RGBA{0, 128, 0, 255}, draw.Src) // Фон. - зелений

	// TODO: Змінити колір фону та додати відображення фігури у вашому варіанті.
	//синя фігура за варіантом Т-270
	blue := color.RGBA{0, 0, 255, 255}
	cx := pw.figureX
	cy := pw.figureY

	//вертикальна частина (ніжка)
	pw.w.Fill(image.Rect(cx-30, cy-100, cx+30, cy+100), blue, draw.Src)

	//горизонтальна частина (капелюшок)
	pw.w.Fill(image.Rect(cx-100, cy-30, cx+30, cy+30), blue, draw.Src)

	// Малювання білої рамки.
	for _, br := range imageutil.Border(pw.sz.Bounds(), 10) {
		pw.w.Fill(br, color.White, draw.Src)
	}
}
