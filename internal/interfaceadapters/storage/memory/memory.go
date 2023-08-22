package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"receipt-processor-challenge/internal/domain/receipt"
)

const (
	get operation = iota
	set
)

var (
	once    sync.Once
	storage map[uuid.UUID]receipt.Points
	engine  Engine

	ErrNotFound = errors.New("points not found")

	errTimeOut   = errors.New("time out")
	errInvalidID = errors.New("invalid id")
)

type operation int

type request struct {
	op  operation
	in  chan payload
	out chan payload
}

type payload struct {
	id   uuid.UUID
	data receipt.Points
	err  error
}

type Engine struct {
	req chan request
}

func New(ctx context.Context, logLevel logrus.Level) *Engine {
	once.Do(func() {
		logrus.SetLevel(logLevel)
		engine.req = make(chan request)
		go engine.start(ctx)
	})

	return &engine
}

func (e *Engine) start(ctx context.Context) {
	storage = make(map[uuid.UUID]receipt.Points)

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-e.req:
			switch req.op {
			case get:
				e.get(req)
			case set:
				e.save(req)
			}
		}
	}
}

func (e *Engine) save(req request) {
	defer func() {
		logrus.Debug("save called")
	}()

	data := <-req.in
	storage[data.id] = data.data

	req.out <- payload{err: nil}
}

func (e *Engine) get(req request) {
	defer func() {
		logrus.Debug("get called")
	}()

	data := <-req.in

	if data.id == uuid.Nil {
		req.out <- payload{
			id:  data.id,
			err: errInvalidID,
		}

		return
	}

	val, ok := storage[data.id]
	if !ok {
		req.out <- payload{
			id:  data.id,
			err: ErrNotFound,
		}

		return
	}

	req.out <- payload{
		id:   data.id,
		data: val,
		err:  nil,
	}
}

func (e *Engine) Save(ctx context.Context, receipt receipt.Points) (uuid.UUID, error) {
	id := uuid.New()

	saveRequest := request{
		op:  set,
		in:  make(chan payload),
		out: make(chan payload),
	}

	select {
	case <-ctx.Done():
		return uuid.Nil, errTimeOut
	case e.req <- saveRequest:
		saveRequest.in <- payload{
			id:   id,
			data: receipt,
		}
	}

	<-saveRequest.out

	return id, nil
}

func (e *Engine) Get(ctx context.Context, id uuid.UUID) (*receipt.Points, error) {
	getRequest := request{
		op:  get,
		in:  make(chan payload),
		out: make(chan payload),
	}

	select {
	case <-ctx.Done():
		return nil, errTimeOut
	case e.req <- getRequest:
		getRequest.in <- payload{id: id}
	}

	data := <-getRequest.out

	return &data.data, data.err
}
