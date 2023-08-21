package drain

import "io"

//go:generate go run github.com/AikoPanel/Xray-core/common/errors/errorgen

type Drainer interface {
	AcknowledgeReceive(size int)
	Drain(reader io.Reader) error
}
