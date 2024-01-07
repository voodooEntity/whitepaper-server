package main

import (
	"fmt"
	"github.com/voodooEntity/archivist"
	"github.com/voodooEntity/gits"
	"github.com/voodooEntity/gits/src/query"
	"github.com/voodooEntity/gits/src/transport"
	"github.com/voodooEntity/gits/src/types"
	"github.com/voodooEntity/whitepaper-server/src/config"
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

	gits.MapTransportData(
		transport.TransportEntity{
			ID:         -1,
			Context:    "ABC",
			Type:       "Something",
			Value:      "asdasd",
			Properties: map[string]string{"foo": "bar"},
		})

	qry := query.New().Read("Something")
	ret := query.Execute(qry)
	fmt.Print("%+v", ret)
}
