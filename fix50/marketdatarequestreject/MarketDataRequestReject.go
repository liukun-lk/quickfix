//Package marketdatarequestreject msg type = Y.
package marketdatarequestreject

import (
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/fix50/mdrjctgrp"
	"github.com/quickfixgo/quickfix/fixt11"
)

//Message is a MarketDataRequestReject FIX Message
type Message struct {
	FIXMsgType string `fix:"Y"`
	fixt11.Header
	//MDReqID is a required field for MarketDataRequestReject.
	MDReqID string `fix:"262"`
	//MDReqRejReason is a non-required field for MarketDataRequestReject.
	MDReqRejReason *string `fix:"281"`
	//MDRjctGrp is a non-required component for MarketDataRequestReject.
	MDRjctGrp *mdrjctgrp.MDRjctGrp
	//Text is a non-required field for MarketDataRequestReject.
	Text *string `fix:"58"`
	//EncodedTextLen is a non-required field for MarketDataRequestReject.
	EncodedTextLen *int `fix:"354"`
	//EncodedText is a non-required field for MarketDataRequestReject.
	EncodedText *string `fix:"355"`
	fixt11.Trailer
}

//Marshal converts Message to a quickfix.Message instance
func (m Message) Marshal() quickfix.Message { return quickfix.Marshal(m) }

//New returns an initialized MarketDataRequestReject instance
func New(mdreqid string) *Message {
	var m Message
	m.SetMDReqID(mdreqid)
	return &m
}

func (m *Message) SetMDReqID(v string)                { m.MDReqID = v }
func (m *Message) SetMDReqRejReason(v string)         { m.MDReqRejReason = &v }
func (m *Message) SetMDRjctGrp(v mdrjctgrp.MDRjctGrp) { m.MDRjctGrp = &v }
func (m *Message) SetText(v string)                   { m.Text = &v }
func (m *Message) SetEncodedTextLen(v int)            { m.EncodedTextLen = &v }
func (m *Message) SetEncodedText(v string)            { m.EncodedText = &v }

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
	return enum.BeginStringFIX50, "Y", r
}
