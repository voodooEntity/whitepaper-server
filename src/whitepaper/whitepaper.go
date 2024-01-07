package whitepaper

import (
	"github.com/voodooEntity/gits"
	"github.com/voodooEntity/gits/src/query"
	"github.com/voodooEntity/gits/src/transport"
	"strconv"
)

var entityName = "Whitepaper"

type WhitePaper struct {
	Hash       string
	Content    string
	ClientID   string
	InstanceId int
}

func (self *WhitePaper) StoreOrUpdate() *WhitePaper {
	gits.MapTransportData(
		transport.TransportEntity{
			ID:         0,
			Type:       entityName,
			Context:    strconv.Itoa(self.InstanceId),
			Value:      "",
			Properties: map[string]string{"Hash": self.Hash, "ClientID": self.ClientID, "Content": self.Content},
		})
	return self
}

func Load(instanceId int) *WhitePaper {
	qry := query.New().Read(entityName).Match("Value", "==", strconv.Itoa(instanceId))
	res := query.Execute(qry)
	if res.Amount == 0 {
		return nil
	}
	return &WhitePaper{
		InstanceId: res.Entities[0].ID,
		Content:    res.Entities[0].Value,
		Hash:       strconv.Itoa(res.Entities[0].Version),
	}
}

func (self *WhitePaper) Update() *WhitePaper {
	query.New().Update(entityName).Match("ID", "==", strconv.Itoa(self.InstanceId)).Set("Value", self.Content)
	return self
}
