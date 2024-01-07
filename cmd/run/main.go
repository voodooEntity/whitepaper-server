package main

import (
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits"
	"github.com/voodooEntity/gits/src/types"
	"github.com/voodooEntity/whitepaper-server/src/config"
	"github.com/voodooEntity/whitepaper-server/src/server"
)

func main() {
	config.Init(make(map[string]string))

	archivist.Init(config.GetValue("LOG_LEVEL"), config.GetValue("LOG_TARGET"), config.GetValue("LOG_PATH"))

	persistence := false
	if "active" == config.GetValue("PERSISTENCE") {
		persistence = true
	}

	gits.Init(types.PersistenceConfig{
		RotationEntriesMax:           1000000,
		Active:                       persistence,
		PersistenceChannelBufferSize: 10000000,
	})

	server.Start()
}
