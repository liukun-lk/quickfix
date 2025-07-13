package gorm

import (
	"time"
)

type Sessions struct {
	BeginString      string    `gorm:"column:beginstring;primaryKey;type:varchar(8)"`
	SenderCompID     string    `gorm:"column:sendercompid;primaryKey;type:varchar(64)"`
	SenderSubID      string    `gorm:"column:sendersubid;primaryKey;type:varchar(64)"`
	SenderLocID      string    `gorm:"column:senderlocid;primaryKey;type:varchar(64)"`
	TargetCompID     string    `gorm:"column:targetcompid;primaryKey;type:varchar(64)"`
	TargetSubID      string    `gorm:"column:targetsubid;primaryKey;type:varchar(64)"`
	TargetLocID      string    `gorm:"column:targetlocid;primaryKey;type:varchar(64)"`
	SessionQualifier string    `gorm:"column:session_qualifier;primaryKey;type:varchar(64)"`
	CreationTime     time.Time `gorm:"column:creation_time"`
	IncomingSeqNum   int       `gorm:"column:incoming_seqnum"`
	OutgoingSeqNum   int       `gorm:"column:outgoing_seqnum"`
}

func (g Sessions) TableName() string {
	return "sessions"
}

type Messages struct {
	BeginString      string `gorm:"column:beginstring;primaryKey;type:varchar(8)"`
	SenderCompID     string `gorm:"column:sendercompid;primaryKey;type:varchar(64)"`
	SenderSubID      string `gorm:"column:sendersubid;primaryKey;type:varchar(64)"`
	SenderLocID      string `gorm:"column:senderlocid;primaryKey;type:varchar(64)"`
	TargetCompID     string `gorm:"column:targetcompid;primaryKey;type:varchar(64)"`
	TargetSubID      string `gorm:"column:targetsubid;primaryKey;type:varchar(64)"`
	TargetLocID      string `gorm:"column:targetlocid;primaryKey;type:varchar(64)"`
	SessionQualifier string `gorm:"column:session_qualifier;primaryKey;type:varchar(64)"`
	Message          string `gorm:"column:message;type:text"`
	MsgSeqNum        int64  `gorm:"column:msgseqnum;primaryKey"`
}

func (g Messages) TableName() string {
	return "messages"
}
