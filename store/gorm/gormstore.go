package gorm

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	gormio "gorm.io/gorm"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/config"
)

type gormStoreFactory struct {
	settings *quickfix.Settings
	db       *gormio.DB
}

func NewGormStoreFactory(settings *quickfix.Settings, db *gormio.DB) quickfix.MessageStoreFactory {
	return gormStoreFactory{settings: settings, db: db}
}

type gromStore struct {
	sessionID quickfix.SessionID
	cache     quickfix.MessageStore
	db        *gormio.DB
}

func (f gormStoreFactory) Create(sessionID quickfix.SessionID) (msgStore quickfix.MessageStore, err error) {
	var dynamicSessions bool
	if f.settings.GlobalSettings().HasSetting(config.DynamicSessions) {
		if dynamicSessions, err = f.settings.GlobalSettings().BoolSetting(config.DynamicSessions); err != nil {
			return
		}
	}
	log.Println("f.settings.SessionSettings()", f.settings.SessionSettings())
	_, ok := f.settings.SessionSettings()[sessionID]
	if !ok && !dynamicSessions {
		return nil, fmt.Errorf("unknown session: %v", sessionID)
	}

	memStore, memErr := quickfix.NewMemoryStoreFactory().Create(sessionID)
	if memErr != nil {
		err = errors.Wrap(memErr, "cache creation")
		return
	}

	store := &gromStore{
		sessionID: sessionID,
		cache:     memStore,
		db:        f.db,
	}
	err = store.initTables()
	if err != nil {
		err = errors.Wrap(err, "initTables err")
		return
	}
	if err = store.cache.Reset(); err != nil {
		err = errors.Wrap(err, "cache reset")
		return
	}
	if err = store.populateCache(); err != nil {
		return nil, err
	}
	return store, nil

}

func (store *gromStore) initTables() (err error) {
	if !store.db.Migrator().HasTable("sessions") {
		err = store.db.Migrator().CreateTable(&GormSessions{})
		if err != nil {
			return errors.Wrap(err, "gromStore.initTables err")
		}
	}
	if !store.db.Migrator().HasTable("messages") {
		err = store.db.Migrator().CreateTable(&GormMessages{})
		if err != nil {
			return errors.Wrap(err, "gromStore.initTables err")
		}
	}
	return nil
}

// Reset deletes the store records and sets the seqnums back to 1.
func (store *gromStore) Reset() error {
	s := store.sessionID
	err := store.db.Exec(`DELETE FROM messages
	WHERE beginstring=? AND session_qualifier=?
	AND sendercompid=? AND sendersubid=? AND senderlocid=?
	AND targetcompid=? AND targetsubid=? AND targetlocid=?`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Error
	if err != nil {
		return err
	}
	if err = store.cache.Reset(); err != nil {
		return err
	}
	err = store.db.Table(`sessions`).Where(`beginstring=? AND session_qualifier=?
	AND sendercompid=? AND sendersubid=? AND senderlocid=?
	AND targetcompid=? AND targetsubid=? AND targetlocid=?`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Updates(map[string]interface{}{
		"creation_time":   store.cache.CreationTime(),
		"incoming_seqnum": store.cache.NextTargetMsgSeqNum(),
		"outgoing_seqnum": store.cache.NextSenderMsgSeqNum(),
	}).Error
	return err
}

// Refresh reloads the store from the database.
func (store *gromStore) Refresh() error {
	if err := store.cache.Reset(); err != nil {
		return err
	}
	return store.populateCache()
}

func (store *gromStore) populateCache() error {
	dest := GormSessions{}
	s := store.sessionID
	err := store.db.Table(`sessions`).Where(`beginstring=? AND session_qualifier=?
	  AND sendercompid=? AND sendersubid=? AND senderlocid=?
	  AND targetcompid=? AND targetsubid=? AND targetlocid=?`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).First(&dest).Error
	if err == nil {
		store.cache.SetCreationTime(dest.CreationTime)
		if err = store.cache.SetNextTargetMsgSeqNum(dest.IncomingSeqNum); err != nil {
			return errors.Wrap(err, "cache set next target")
		}
		if err = store.cache.SetNextSenderMsgSeqNum(dest.OutgoingSeqNum); err != nil {
			return errors.Wrap(err, "cache set next sender")
		}
		return nil
	}
	if err == gormio.ErrRecordNotFound {
		return store.db.Exec(`INSERT INTO sessions (
			creation_time, incoming_seqnum, outgoing_seqnum,
			beginstring, session_qualifier,
			sendercompid, sendersubid, senderlocid,
			targetcompid, targetsubid, targetlocid)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, store.cache.CreationTime(),
			store.cache.NextTargetMsgSeqNum(),
			store.cache.NextSenderMsgSeqNum(),
			s.BeginString, s.Qualifier,
			s.SenderCompID, s.SenderSubID, s.SenderLocationID,
			s.TargetCompID, s.TargetSubID, s.TargetLocationID).Error
	}
	return err
}

// NextSenderMsgSeqNum returns the next MsgSeqNum that will be sent.
func (store *gromStore) NextSenderMsgSeqNum() int {
	return store.cache.NextSenderMsgSeqNum()
}

// NextTargetMsgSeqNum returns the next MsgSeqNum that should be received.
func (store *gromStore) NextTargetMsgSeqNum() int {
	return store.cache.NextTargetMsgSeqNum()
}

// SetNextSenderMsgSeqNum sets the next MsgSeqNum that will be sent.
func (store *gromStore) SetNextSenderMsgSeqNum(next int) error {
	s := store.sessionID

	err := store.db.Table(`sessions`).Where(`beginstring=? AND session_qualifier=?
	AND sendercompid=? AND sendersubid=? AND senderlocid=?
	AND targetcompid=? AND targetsubid=? AND targetlocid=?`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Update(`outgoing_seqnum`, next).Error
	if err != nil {
		return err
	}
	return store.cache.SetNextSenderMsgSeqNum(next)
}

// SetNextTargetMsgSeqNum sets the next MsgSeqNum that should be received.
func (store *gromStore) SetNextTargetMsgSeqNum(next int) error {
	s := store.sessionID

	err := store.db.Table(`sessions`).Where(`beginstring=? AND session_qualifier=?
	AND sendercompid=? AND sendersubid=? AND senderlocid=?
	AND targetcompid=? AND targetsubid=? AND targetlocid=?`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Update(`incoming_seqnum`, next).Error
	if err != nil {
		return err
	}
	return store.cache.SetNextTargetMsgSeqNum(next)
}

// IncrNextSenderMsgSeqNum increments the next MsgSeqNum that will be sent.
func (store *gromStore) IncrNextSenderMsgSeqNum() error {
	if err := store.cache.IncrNextSenderMsgSeqNum(); err != nil {
		return errors.Wrap(err, "cache incr next")
	}
	return store.SetNextSenderMsgSeqNum(store.cache.NextSenderMsgSeqNum())
}

// IncrNextTargetMsgSeqNum increments the next MsgSeqNum that should be received.
func (store *gromStore) IncrNextTargetMsgSeqNum() error {
	if err := store.cache.IncrNextTargetMsgSeqNum(); err != nil {
		return errors.Wrap(err, "cache incr next")
	}
	return store.SetNextTargetMsgSeqNum(store.cache.NextTargetMsgSeqNum())
}

// CreationTime returns the creation time of the store.
func (store *gromStore) CreationTime() time.Time {
	return store.cache.CreationTime()
}

// SetCreationTime is a no-op for GormStore.
func (store *gromStore) SetCreationTime(_ time.Time) {
}

func (store *gromStore) SaveMessage(seqNum int, msg []byte) error {
	s := store.sessionID
	err := store.db.Exec(`INSERT INTO messages (
		msgseqnum, message,
		beginstring, session_qualifier,
		sendercompid, sendersubid, senderlocid,
		targetcompid, targetsubid, targetlocid)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, seqNum, string(msg),
		s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Error
	if err != nil {
		var counter int64
		store.db.Table("messages").Where(`beginstring=? AND session_qualifier=?
		AND sendercompid=? AND sendersubid=? AND senderlocid=?
		AND targetcompid=? AND targetsubid=? AND targetlocid=?
		AND msgseqnum=?`, s.BeginString, s.Qualifier,
			s.SenderCompID, s.SenderSubID, s.SenderLocationID,
			s.TargetCompID, s.TargetSubID, s.TargetLocationID,
			seqNum).Limit(1).Count(&counter)
		// If it is determined that the message is repeated, skip this insertion
		if counter == 1 {
			return nil
		}
	}
	return err
}

func (store *gromStore) SaveMessageAndIncrNextSenderMsgSeqNum(seqNum int, msg []byte) error {
	s := store.sessionID
	tx := store.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer tx.Rollback()

	err := tx.Exec(`INSERT INTO messages (
		msgseqnum, message,
		beginstring, session_qualifier,
		sendercompid, sendersubid, senderlocid,
		targetcompid, targetsubid, targetlocid)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, seqNum, string(msg),
		s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Error
	if err != nil {
		return err
	}

	next := store.cache.NextSenderMsgSeqNum() + 1
	err = tx.Exec(`UPDATE sessions SET outgoing_seqnum = ?
		WHERE beginstring = ? AND session_qualifier = ?
		AND sendercompid = ? AND sendersubid = ? AND senderlocid = ?
		AND targetcompid = ? AND targetsubid = ? AND targetlocid = ?`, next,
		s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID).Error
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return store.cache.SetNextSenderMsgSeqNum(next)
}

func (store *gromStore) IterateMessages(beginSeqNum, endSeqNum int, cb func([]byte) error) error {
	s := store.sessionID
	rows, err := store.db.Raw(`SELECT message FROM messages
	WHERE beginstring=? AND session_qualifier=?
	AND sendercompid=? AND sendersubid=? AND senderlocid=?
	AND targetcompid=? AND targetsubid=? AND targetlocid=?
	AND msgseqnum>=? AND msgseqnum<=?
	ORDER BY msgseqnum`, s.BeginString, s.Qualifier,
		s.SenderCompID, s.SenderSubID, s.SenderLocationID,
		s.TargetCompID, s.TargetSubID, s.TargetLocationID,
		beginSeqNum, endSeqNum).Rows()
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var message string
		if err = rows.Scan(&message); err != nil {
			return err
		} else if err = cb([]byte(message)); err != nil {
			return err
		}
	}

	return rows.Err()
}

func (store *gromStore) GetMessages(beginSeqNum, endSeqNum int) ([][]byte, error) {
	var msgs [][]byte
	err := store.IterateMessages(beginSeqNum, endSeqNum, func(msg []byte) error {
		msgs = append(msgs, msg)
		return nil
	})
	return msgs, err
}

// Close closes the store's database connection
func (store *gromStore) Close() error {
	if store.db != nil {
		db, err := store.db.DB()
		if err != nil {
			db.Close()
		}
		store.db = nil
	}
	return nil
}
