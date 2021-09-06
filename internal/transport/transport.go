package transport

import (
	"context"
	"github.com/belliel/multiplexer/internal/transport/http"
)

const (
	defaultAddr  = ":8080"
	defaultDebug = true
)

type TransportType int

const (
	HTTP TransportType = iota
	// other types of transport
)

type Transporter interface {
	Listen() error
	ListenCtxForGracefulShutdown()
	Shutdown(ctx context.Context) error
	WaitForGracefulShutdown()
}

type Transport struct {
	masterCtx context.Context

	transporter *Transporter
}

type Builder struct {
	mainCtx       context.Context
	transportType TransportType
	transport     *Transporter
	addr          string
	debug         bool
}

func NewTransportBuilder(ctx context.Context, transportType TransportType) *Builder {
	return &Builder{
		mainCtx:       ctx,
		transportType: transportType,
		addr:          defaultAddr,
		debug:         defaultDebug,
	}
}

func (b *Builder) WithAddr(addr string) *Builder {
	b.addr = addr
	return b
}

func (b *Builder) WithDebug(debug bool) *Builder {
	b.debug = debug
	return b
}

func (b *Builder) Build() Transporter {
	switch b.transportType {
	case HTTP:
		return http.NewServer(b.mainCtx, b.debug, b.addr)
	default:
		panic("Transport type is not implemented")
	}
}
