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
	qry := query.New().Read("Whitepaper").Match("Value", "==", self.InstanceId)
	ret := query.Execute(qry)
	if 0 == ret.Amount {
		gits.MapTransportData(
			transport.TransportEntity{
				ID:         0,
				Type:       entityName,
				Context:    "",
				Value:      self.InstanceId,
				Properties: map[string]string{"Hash": self.Hash, "ClientID": self.ClientID, "Content": self.Content},
			})
	} else {
		qry = query.New().Update("Whitepaper").Match("Value", "==", self.InstanceId).Set("Properties.Content", self.Content).Set("Properties.Hash", self.Hash).Set("Properties.ClientID", self.ClientID)
		query.Execute(qry)
	}

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
