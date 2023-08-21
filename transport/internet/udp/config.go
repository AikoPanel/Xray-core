package udp

import (
	"github.com/AikoPanel/Xray-core/common"
	"github.com/AikoPanel/Xray-core/transport/internet"
)

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}
