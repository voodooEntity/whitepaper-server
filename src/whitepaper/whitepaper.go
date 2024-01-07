package whitepaper

import (
	"github.com/voodooEntity/gits"
	"github.com/voodooEntity/gits/src/query"
	"github.com/voodooEntity/gits/src/transport"
)

var entityName = "Whitepaper"

type WhitePaper struct {
	Hash       string
	Content    string
	ClientID   string
	InstanceId string
}

func (self *WhitePaper) StoreOrUpdate() *WhitePaper {
	gits.MapTransportData(
		transport.TransportEntity{
			ID:         0,
			Type:       entityName,
			Context:    self.InstanceId,
			Value:      "",
			Properties: map[string]string{"Hash": self.Hash, "ClientID": self.ClientID, "Content": self.Content},
		})
	return self
}

func Load(instanceId string) *WhitePaper {
	qry := query.New().Read(entityName).Match("Value", "==", instanceId)
	res := query.Execute(qry)
	if res.Amount == 0 {
		return nil
	}
	return &WhitePaper{
		InstanceId: res.Entities[0].Value,
		Content:    res.Entities[0].Properties["Content"],
		Hash:       res.Entities[0].Properties["Hash"],
	}
}
