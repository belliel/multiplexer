package transport

import (
	"context"
	"strconv"
)

type TransportType int

const (
	HTTP TransportType = iota
	// other types of transport
)

type Transporter interface {
	Listen()
	Shutdown()
}

type Transport struct {
	masterCtx context.Context

	transporter *Transporter
}

type Builder struct {
	mainCtx       context.Context
	transportType TransportType
	transport     *Transporter
	port          string
	debug         bool
}

func NewTransportBuilder(ctx context.Context, transportType TransportType) *Builder {
	return &Builder{
		mainCtx:       ctx,
		transportType: transportType,
		port:          "3000",
		debug:         true,
	}
}

func (b *Builder) WithPort(port string) *Builder {
	b.port = port
	return b
}

func (b *Builder) WithDebug(debug bool) *Builder {
	b.debug = debug
	return b
}

func (b *Builder) Build() *Transporter {
	switch b.transportType {
	case HTTP:
		return NewHttpServer{}
	}
}
