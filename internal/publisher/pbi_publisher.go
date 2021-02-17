package publisher

import (
	"bytes"
	"go.uber.org/zap"
	"net/http"
)

type PBIPublisher struct {
	Url string
	Logger *zap.Logger
}

func (p *PBIPublisher) Publish(data []byte) {
	_, err := http.Post(p.Url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		p.Logger.Error("Publish error", zap.Error(err))
	}
}