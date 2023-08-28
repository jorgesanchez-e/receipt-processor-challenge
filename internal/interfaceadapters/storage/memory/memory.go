package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"receipt-processor-challenge/internal/domain/receipt"

	"github.com/google/uuid"
)

const (
	get operation = iota
	set

	defaultTimeOut = time.Second
)

var (
	once    sync.Once                    //nolint:gochecknoglobals
	storage map[uuid.UUID]receipt.Points //nolint:gochecknoglobals
	engine  Engine                       //nolint:gochecknoglobals

	ErrNotFound = errors.New("points not found")

	errTimeOut   = errors.New("time out")
	errInvalidID = errors.New("invalid id")
)

type operation int

type request struct {
	ctx context.Context
	op  operation
	in  chan payload
	out chan payload
}

type payload struct {
	id   uuid.UUID
	data *receipt.Points
	err  error
}

type Engine struct {
	req       chan request
	opTimeOut time.Duration
}

func New(ctx context.Context) *Engine {
	once.Do(func() {
		engine.req = make(chan request)
		engine.opTimeOut = defaultTimeOut
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
	var data payload

	select {
	case data = <-req.in:
	case <-req.ctx.Done():
		return
	}

	storage[data.id] = *data.data

	defer close(req.out)

	pload := payload{
		id:   uuid.Nil,
		data: nil,
		err:  nil,
	}

	select {
	case req.out <- pload:
	case <-req.ctx.Done():
	}
}

func (e *Engine) get(req request) {
	var data payload

	select {
	case data = <-req.in:
	case <-req.ctx.Done():
		return
	}

	var pload *payload

	defer close(req.out)
	defer func() {
		select {
		case <-req.ctx.Done():
		case req.out <- *pload:
		}
	}()

	if data.id == uuid.Nil {
		pload = &payload{
			id:   data.id,
			data: nil,
			err:  errInvalidID,
		}

		return
	}

	val, ok := storage[data.id]
	if !ok {
		pload = &payload{
			id:   data.id,
			data: nil,
			err:  ErrNotFound,
		}

		return
	}

	pload = &payload{
		id:   data.id,
		data: &val,
		err:  nil,
	}
}

func (e *Engine) Save(ctx context.Context, receipt receipt.Points) (uuid.UUID, error) {
	if _, deadLineSet := ctx.Deadline(); !deadLineSet {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, e.opTimeOut)
		defer cancel()
	}

	newID := uuid.New()

	saveRequest := request{
		ctx: ctx,
		op:  set,
		in:  make(chan payload),
		out: make(chan payload),
	}

	pload := payload{
		id:   newID,
		data: &receipt,
		err:  nil,
	}

	select {
	case <-ctx.Done():
		return uuid.Nil, ctx.Err()
	case e.req <- saveRequest:
		select {
		case saveRequest.in <- pload:
			close(saveRequest.in)
		case <-ctx.Done():
			return uuid.Nil, ctx.Err()
		}
	}

	<-saveRequest.out

	return newID, nil
}

func (e *Engine) Get(ctx context.Context, uid uuid.UUID) (*receipt.Points, error) {
	if _, deadLineSet := ctx.Deadline(); !deadLineSet {
		var cancel context.CancelFunc

		ctx, cancel = context.WithTimeout(ctx, e.opTimeOut)
		defer cancel()
	}

	getRequest := request{
		ctx: ctx,
		op:  get,
		in:  make(chan payload),
		out: make(chan payload),
	}

	pload := payload{
		id:   uid,
		data: nil,
		err:  nil,
	}

	select {
	case <-ctx.Done():
		return nil, errTimeOut
	case e.req <- getRequest:
		select {
		case getRequest.in <- pload:
			close(getRequest.in)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	data := <-getRequest.out

	return data.data, data.err
}
