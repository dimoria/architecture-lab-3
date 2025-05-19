package painter

import (
	"image"
	"image/color"
	"sync"

	"golang.org/x/exp/shiny/screen"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	next     screen.Texture
	prev     screen.Texture
	mq       *messageQueue
	stopChan chan struct{}
	state    *State
}

type State struct {
	bgColor   color.Color
	bgRect    *image.Rectangle
	figures   []image.Point
	moveDelta image.Point
}

func NewLoop() *Loop {
	return &Loop{
		mq:       newMessageQueue(),
		stopChan: make(chan struct{}),
		state: &State{
			bgColor: color.Black,
		},
	}
}

func (l *Loop) Start(s screen.Screen) {
	var err error
	l.next, err = s.NewTexture(image.Point{800, 800})
	if err != nil {
		panic(err)
	}
	l.prev, err = s.NewTexture(image.Point{800, 800})
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			op := l.mq.pull()
			if op == nil {
				close(l.stopChan)
				return
			}

			if update := op.Do(l.next, l.state); update {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
				l.next.Fill(l.next.Bounds(), l.state.bgColor, screen.Src)
			}
		}
	}()
}

func (l *Loop) Post(op Operation) {
	l.mq.push(op)
}

func (l *Loop) StopAndWait() {
	l.mq.shutdown()
	<-l.stopChan
}

type messageQueue struct {
	ops     []Operation
	mu      sync.Mutex
	cond    *sync.Cond
	blocked bool
}

func newMessageQueue() *messageQueue {
	mq := &messageQueue{}
	mq.cond = sync.NewCond(&mq.mu)
	return mq
}

func (mq *messageQueue) push(op Operation) {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.ops = append(mq.ops, op)
	mq.cond.Signal()
}

func (mq *messageQueue) pull() Operation {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	for len(mq.ops) == 0 && !mq.blocked {
		mq.cond.Wait()
	}

	if len(mq.ops) == 0 {
		return nil
	}

	op := mq.ops[0]
	mq.ops = mq.ops[1:]
	return op
}

func (mq *messageQueue) shutdown() {
	mq.mu.Lock()
	defer mq.mu.Unlock()

	mq.blocked = true
	mq.cond.Broadcast()
}
