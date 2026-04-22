package lib_ws

import (
	"net/http"

	"github.com/pkg/errors"
)

type httpServerImpl struct {
	*baseServer
}

func NewHttpServer(address string, options ...Option) (Server, error) {
	return &httpServerImpl{
		baseServer: newBaseServer(address, options...),
	}, nil
}

func (w httpServerImpl) Start() error {
	if err := w.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
