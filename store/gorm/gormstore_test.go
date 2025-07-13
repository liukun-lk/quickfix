package gorm

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/internal/testsuite"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	gormio "gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormStoreTestSuite runs all tests in the MessageStoreTestSuite against the GormStore implementation.
type GormStoreTestSuite struct {
	testsuite.StoreTestSuite
	db *gormio.DB
}

func (suite *GormStoreTestSuite) SetupTest() {
	var err error
	dsn := "host=127.0.0.1 user=postgres dbname=lb_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	suite.db, err = gormio.Open(postgres.Open(dsn), &gormio.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.Nil(suite.T(), err)

	// 创建表
	err = suite.db.AutoMigrate(&Sessions{}, &Messages{})
	require.Nil(suite.T(), err)

	// 创建 store
	sessionID := quickfix.SessionID{BeginString: "FIX.4.4", SenderCompID: "SENDER", TargetCompID: "TARGET"}
	settings, err := quickfix.ParseSettings(strings.NewReader(fmt.Sprintf(`
[SESSION]
BeginString=%s
SenderCompID=%s
TargetCompID=%s`, sessionID.BeginString, sessionID.SenderCompID, sessionID.TargetCompID)))
	require.Nil(suite.T(), err)

	factory := NewGormStoreFactory(settings, suite.db)
	suite.MsgStore, err = factory.Create(sessionID)
	require.Nil(suite.T(), err)
}

func (suite *GormStoreTestSuite) TearDownTest() {
	if suite.db != nil {
		sqlDB, err := suite.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if suite.MsgStore != nil {
		suite.MsgStore.Close()
	}

	// 删除测试数据库文件
	os.Remove("test.db")
}

func TestGormStoreTestSuite(t *testing.T) {
	suite.Run(t, new(GormStoreTestSuite))
}

// 额外的特定于 GormStore 的测试。
func TestGormStoreSpecificFeatures(t *testing.T) {
	dsn := "host=127.0.0.1 user=postgres dbname=lb_test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gormio.Open(postgres.Open(dsn), &gormio.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// 创建表
	err = db.AutoMigrate(&Sessions{}, &Messages{})
	require.NoError(t, err)

	sessionID := quickfix.SessionID{BeginString: "FIX.4.2", TargetCompID: "IB", SenderCompID: "LB"}
	settings, err := quickfix.ParseSettings(strings.NewReader(fmt.Sprintf(`
[SESSION]
BeginString=%s
SenderCompID=%s
TargetCompID=%s`, sessionID.BeginString, sessionID.SenderCompID, sessionID.TargetCompID)))
	require.NoError(t, err)

	factory := NewGormStoreFactory(settings, db)
	store, err := factory.Create(sessionID)
	require.NoError(t, err)
	defer store.Close()

	// 测试保存重复消息
	err = store.SaveMessage(1, []byte(`test msg`))
	require.NoError(t, err)
	err = store.SaveMessage(1, []byte(`test msg updated`))
	require.NoError(t, err)

	// 测试设置和获取序列号
	err = store.SetNextSenderMsgSeqNum(10)
	require.NoError(t, err)
	require.Equal(t, 10, store.NextSenderMsgSeqNum())
}
