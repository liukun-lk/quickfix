//Package ordercancelreplacerequest msg type = G.
package ordercancelreplacerequest

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/fix40"
	"time"
)

//Message is a OrderCancelReplaceRequest FIX Message
type Message struct {
	FIXMsgType string `fix:"G"`
	fix40.Header
	//OrderID is a non-required field for OrderCancelReplaceRequest.
	OrderID *string `fix:"37"`
	//ClientID is a non-required field for OrderCancelReplaceRequest.
	ClientID *string `fix:"109"`
	//ExecBroker is a non-required field for OrderCancelReplaceRequest.
	ExecBroker *string `fix:"76"`
	//OrigClOrdID is a required field for OrderCancelReplaceRequest.
	OrigClOrdID string `fix:"41"`
	//ClOrdID is a required field for OrderCancelReplaceRequest.
	ClOrdID string `fix:"11"`
	//ListID is a non-required field for OrderCancelReplaceRequest.
	ListID *string `fix:"66"`
	//Account is a non-required field for OrderCancelReplaceRequest.
	Account *string `fix:"1"`
	//SettlmntTyp is a non-required field for OrderCancelReplaceRequest.
	SettlmntTyp *string `fix:"63"`
	//FutSettDate is a non-required field for OrderCancelReplaceRequest.
	FutSettDate *string `fix:"64"`
	//HandlInst is a required field for OrderCancelReplaceRequest.
	HandlInst string `fix:"21"`
	//ExecInst is a non-required field for OrderCancelReplaceRequest.
	ExecInst *string `fix:"18"`
	//MinQty is a non-required field for OrderCancelReplaceRequest.
	MinQty *int `fix:"110"`
	//MaxFloor is a non-required field for OrderCancelReplaceRequest.
	MaxFloor *int `fix:"111"`
	//ExDestination is a non-required field for OrderCancelReplaceRequest.
	ExDestination *string `fix:"100"`
	//Symbol is a required field for OrderCancelReplaceRequest.
	Symbol string `fix:"55"`
	//SymbolSfx is a non-required field for OrderCancelReplaceRequest.
	SymbolSfx *string `fix:"65"`
	//SecurityID is a non-required field for OrderCancelReplaceRequest.
	SecurityID *string `fix:"48"`
	//IDSource is a non-required field for OrderCancelReplaceRequest.
	IDSource *string `fix:"22"`
	//Issuer is a non-required field for OrderCancelReplaceRequest.
	Issuer *string `fix:"106"`
	//SecurityDesc is a non-required field for OrderCancelReplaceRequest.
	SecurityDesc *string `fix:"107"`
	//Side is a required field for OrderCancelReplaceRequest.
	Side string `fix:"54"`
	//OrderQty is a required field for OrderCancelReplaceRequest.
	OrderQty int `fix:"38"`
	//OrdType is a required field for OrderCancelReplaceRequest.
	OrdType string `fix:"40"`
	//Price is a non-required field for OrderCancelReplaceRequest.
	Price *float64 `fix:"44"`
	//StopPx is a non-required field for OrderCancelReplaceRequest.
	StopPx *float64 `fix:"99"`
	//Currency is a non-required field for OrderCancelReplaceRequest.
	Currency *string `fix:"15"`
	//TimeInForce is a non-required field for OrderCancelReplaceRequest.
	TimeInForce *string `fix:"59"`
	//ExpireTime is a non-required field for OrderCancelReplaceRequest.
	ExpireTime *time.Time `fix:"126"`
	//Commission is a non-required field for OrderCancelReplaceRequest.
	Commission *float64 `fix:"12"`
	//CommType is a non-required field for OrderCancelReplaceRequest.
	CommType *string `fix:"13"`
	//Rule80A is a non-required field for OrderCancelReplaceRequest.
	Rule80A *string `fix:"47"`
	//ForexReq is a non-required field for OrderCancelReplaceRequest.
	ForexReq *string `fix:"121"`
	//SettlCurrency is a non-required field for OrderCancelReplaceRequest.
	SettlCurrency *string `fix:"120"`
	//Text is a non-required field for OrderCancelReplaceRequest.
	Text *string `fix:"58"`
	fix40.Trailer
}

//Marshal converts Message to a quickfix.Message instance
func (m Message) Marshal() quickfix.Message { return quickfix.Marshal(m) }

//New returns an initialized OrderCancelReplaceRequest instance
func New(origclordid string, clordid string, handlinst string, symbol string, side string, orderqty int, ordtype string) *Message {
	var m Message
	m.SetOrigClOrdID(origclordid)
	m.SetClOrdID(clordid)
	m.SetHandlInst(handlinst)
	m.SetSymbol(symbol)
	m.SetSide(side)
	m.SetOrderQty(orderqty)
	m.SetOrdType(ordtype)
	return &m
}

func (m *Message) SetOrderID(v string)       { m.OrderID = &v }
func (m *Message) SetClientID(v string)      { m.ClientID = &v }
func (m *Message) SetExecBroker(v string)    { m.ExecBroker = &v }
func (m *Message) SetOrigClOrdID(v string)   { m.OrigClOrdID = v }
func (m *Message) SetClOrdID(v string)       { m.ClOrdID = v }
func (m *Message) SetListID(v string)        { m.ListID = &v }
func (m *Message) SetAccount(v string)       { m.Account = &v }
func (m *Message) SetSettlmntTyp(v string)   { m.SettlmntTyp = &v }
func (m *Message) SetFutSettDate(v string)   { m.FutSettDate = &v }
func (m *Message) SetHandlInst(v string)     { m.HandlInst = v }
func (m *Message) SetExecInst(v string)      { m.ExecInst = &v }
func (m *Message) SetMinQty(v int)           { m.MinQty = &v }
func (m *Message) SetMaxFloor(v int)         { m.MaxFloor = &v }
func (m *Message) SetExDestination(v string) { m.ExDestination = &v }
func (m *Message) SetSymbol(v string)        { m.Symbol = v }
func (m *Message) SetSymbolSfx(v string)     { m.SymbolSfx = &v }
func (m *Message) SetSecurityID(v string)    { m.SecurityID = &v }
func (m *Message) SetIDSource(v string)      { m.IDSource = &v }
func (m *Message) SetIssuer(v string)        { m.Issuer = &v }
func (m *Message) SetSecurityDesc(v string)  { m.SecurityDesc = &v }
func (m *Message) SetSide(v string)          { m.Side = v }
func (m *Message) SetOrderQty(v int)         { m.OrderQty = v }
func (m *Message) SetOrdType(v string)       { m.OrdType = v }
func (m *Message) SetPrice(v float64)        { m.Price = &v }
func (m *Message) SetStopPx(v float64)       { m.StopPx = &v }
func (m *Message) SetCurrency(v string)      { m.Currency = &v }
func (m *Message) SetTimeInForce(v string)   { m.TimeInForce = &v }
func (m *Message) SetExpireTime(v time.Time) { m.ExpireTime = &v }
func (m *Message) SetCommission(v float64)   { m.Commission = &v }
func (m *Message) SetCommType(v string)      { m.CommType = &v }
func (m *Message) SetRule80A(v string)       { m.Rule80A = &v }
func (m *Message) SetForexReq(v string)      { m.ForexReq = &v }
func (m *Message) SetSettlCurrency(v string) { m.SettlCurrency = &v }
func (m *Message) SetText(v string)          { m.Text = &v }

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
	return enum.BeginStringFIX40, "G", r
}
