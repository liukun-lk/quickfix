//Package settlementinstructionrequest msg type = AV.
package settlementinstructionrequest

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/fix44"
	"github.com/quickfixgo/quickfix/fix44/parties"
	"time"
)

//Message is a SettlementInstructionRequest FIX Message
type Message struct {
	FIXMsgType string `fix:"AV"`
	fix44.Header
	//SettlInstReqID is a required field for SettlementInstructionRequest.
	SettlInstReqID string `fix:"791"`
	//TransactTime is a required field for SettlementInstructionRequest.
	TransactTime time.Time `fix:"60"`
	//Parties is a non-required component for SettlementInstructionRequest.
	Parties *parties.Parties
	//AllocAccount is a non-required field for SettlementInstructionRequest.
	AllocAccount *string `fix:"79"`
	//AllocAcctIDSource is a non-required field for SettlementInstructionRequest.
	AllocAcctIDSource *int `fix:"661"`
	//Side is a non-required field for SettlementInstructionRequest.
	Side *string `fix:"54"`
	//Product is a non-required field for SettlementInstructionRequest.
	Product *int `fix:"460"`
	//SecurityType is a non-required field for SettlementInstructionRequest.
	SecurityType *string `fix:"167"`
	//CFICode is a non-required field for SettlementInstructionRequest.
	CFICode *string `fix:"461"`
	//EffectiveTime is a non-required field for SettlementInstructionRequest.
	EffectiveTime *time.Time `fix:"168"`
	//ExpireTime is a non-required field for SettlementInstructionRequest.
	ExpireTime *time.Time `fix:"126"`
	//LastUpdateTime is a non-required field for SettlementInstructionRequest.
	LastUpdateTime *time.Time `fix:"779"`
	//StandInstDbType is a non-required field for SettlementInstructionRequest.
	StandInstDbType *int `fix:"169"`
	//StandInstDbName is a non-required field for SettlementInstructionRequest.
	StandInstDbName *string `fix:"170"`
	//StandInstDbID is a non-required field for SettlementInstructionRequest.
	StandInstDbID *string `fix:"171"`
	fix44.Trailer
}

//Marshal converts Message to a quickfix.Message instance
func (m Message) Marshal() quickfix.Message { return quickfix.Marshal(m) }

//New returns an initialized SettlementInstructionRequest instance
func New(settlinstreqid string, transacttime time.Time) *Message {
	var m Message
	m.SetSettlInstReqID(settlinstreqid)
	m.SetTransactTime(transacttime)
	return &m
}

func (m *Message) SetSettlInstReqID(v string)    { m.SettlInstReqID = v }
func (m *Message) SetTransactTime(v time.Time)   { m.TransactTime = v }
func (m *Message) SetParties(v parties.Parties)  { m.Parties = &v }
func (m *Message) SetAllocAccount(v string)      { m.AllocAccount = &v }
func (m *Message) SetAllocAcctIDSource(v int)    { m.AllocAcctIDSource = &v }
func (m *Message) SetSide(v string)              { m.Side = &v }
func (m *Message) SetProduct(v int)              { m.Product = &v }
func (m *Message) SetSecurityType(v string)      { m.SecurityType = &v }
func (m *Message) SetCFICode(v string)           { m.CFICode = &v }
func (m *Message) SetEffectiveTime(v time.Time)  { m.EffectiveTime = &v }
func (m *Message) SetExpireTime(v time.Time)     { m.ExpireTime = &v }
func (m *Message) SetLastUpdateTime(v time.Time) { m.LastUpdateTime = &v }
func (m *Message) SetStandInstDbType(v int)      { m.StandInstDbType = &v }
func (m *Message) SetStandInstDbName(v string)   { m.StandInstDbName = &v }
func (m *Message) SetStandInstDbID(v string)     { m.StandInstDbID = &v }

//A RouteOut is the callback type that should be implemented for routing Message
type RouteOut func(msg Message, sessionID quickfix.SessionID) quickfix.MessageRejectError

//Route returns the beginstring, message type, and MessageRoute for this Message type
func Route(router RouteOut) (string, string, quickfix.MessageRoute) {
	r := func(msg quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
		m := new(Message)
		if err := quickfix.Unmarshal(msg, m); err != nil {
			return err
		}
		return router(*m, sessionID)
	}
	return enum.BeginStringFIX44, "AV", r
}
