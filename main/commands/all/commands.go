package all

import (
	"github.com/AikoPanel/Xray-core/main/commands/all/api"
	"github.com/AikoPanel/Xray-core/main/commands/all/tls"
	"github.com/AikoPanel/Xray-core/main/commands/base"
)

// go:generate go run github.com/AikoPanel/Xray-core/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		// cmdConvert,
		tls.CmdTLS,
		cmdUUID,
		cmdX25519,
	)
}
