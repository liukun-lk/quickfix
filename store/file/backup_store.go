package file

import (
	"fmt"

	"github.com/quickfixgo/quickfix"
)

const (
	OperationSetNextSenderMsgSeqNum int = iota + 1
	OperationSetNextTargetMsgSeqNum
	OperationSaveMessage
	OperationReset
)

type BackupMessage struct {
	Operation int
	SeqNum    int
	Msg       []byte
}

type backupStoreFactory struct {
	messagesQueue chan *BackupMessage
	backupFactory quickfix.MessageStoreFactory
}

type backupStore struct {
	messagesQueue chan *BackupMessage
	store         quickfix.MessageStore
}

func NewBackupStoreFactory(messagesQueue chan *BackupMessage, backupFactory quickfix.MessageStoreFactory) *backupStoreFactory {
	return &backupStoreFactory{messagesQueue: messagesQueue, backupFactory: backupFactory}
}

func (f backupStoreFactory) Create(sessionID quickfix.SessionID) (msgStore *backupStore, err error) {
	backupStore, err := f.backupFactory.Create(sessionID)
	if err != nil {
		return nil, err
	}

	return newBackupStore(backupStore, f.messagesQueue), nil
}

func newBackupStore(store quickfix.MessageStore, messagesQueue chan *BackupMessage) *backupStore {
	backup := &backupStore{messagesQueue: messagesQueue, store: store}

	backup.start()

	return backup
}

func (s *backupStore) start() {
	if s == nil {
		return
	}

	go func() {
		for message := range s.messagesQueue {
			switch message.Operation {
			case OperationSetNextSenderMsgSeqNum:
				if err := s.store.SetNextSenderMsgSeqNum(message.SeqNum); err != nil {
					fmt.Printf("backup store: SetNextSenderMsgSeqNum error(%v)\n", err)
				}
			case OperationSetNextTargetMsgSeqNum:
				if err := s.store.SetNextTargetMsgSeqNum(message.SeqNum); err != nil {
					fmt.Printf("backup store: SetNextTargetMsgSeqNum error(%v)\n", err)
				}
			case OperationSaveMessage:
				if err := s.store.SaveMessage(message.SeqNum, message.Msg); err != nil {
					fmt.Printf("backup store: SaveMessage error(%v)\n", err)
				}
			case OperationReset:
				if err := s.store.Reset(); err != nil {
					fmt.Printf("backup store: Reset error(%v)\n", err)
				}
			default:
				fmt.Printf("backup store: unsupported operation(%v)\n", message.Operation)
			}
		}
	}()
}

func (s *backupStore) SetNextSenderMsgSeqNum(next int) {
	if s == nil {
		return
	}

	select {
	case s.messagesQueue <- &BackupMessage{Operation: OperationSetNextSenderMsgSeqNum, SeqNum: next}:
	default:
		fmt.Println("encountering a large amount of traffic, drop the SetNextSenderMsgSeqNum operation")
	}
}

func (s *backupStore) SetNextTargetMsgSeqNum(next int) {
	if s == nil {
		return
	}

	select {
	case s.messagesQueue <- &BackupMessage{Operation: OperationSetNextTargetMsgSeqNum, SeqNum: next}:
	default:
		fmt.Println("encountering a large amount of traffic, drop the SetNextTargetMsgSeqNum operation")
	}
}

func (s *backupStore) SaveMessage(seqNum int, msg []byte) {
	if s == nil {
		return
	}

	select {
	case s.messagesQueue <- &BackupMessage{Operation: OperationSaveMessage, SeqNum: seqNum, Msg: msg}:
	default:
		fmt.Println("encountering a large amount of traffic, drop the SaveMessage operation")
	}
}

func (s *backupStore) Reset() {
	if s == nil {
		return
	}

	select {
	case s.messagesQueue <- &BackupMessage{Operation: OperationReset}:
	default:
		fmt.Println("encountering a large amount of traffic, drop the Reset operation")
	}
}
