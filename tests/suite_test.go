package tests

import (
	"github.com/stretchr/testify/suite"
	"main/global"
	"main/middleware/dao"
	"main/model/db"
	"os"
	"testing"
)

type SuiteTest struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	os.Setenv("ENV", "test")
	defer os.Unsetenv("ENV")

	suite.Run(t, new(SuiteTest))
}

func getModels() []interface{} {
	return []interface{}{
		&db.MPCWallet{},
	}
}

// Setup suite
func (t *SuiteTest) SetupSuite() {
	println("[SetupSuite] Setup for all tests")
	global.InitAll()

	// Migrate Table
	for _, val := range getModels() {
		dao.GetDbEngine().AutoMigrate(val)
	}

	// start gin server
}

// Run After All Test Done
func (t *SuiteTest) TearDownSuite() {
	println("[TearDownSuite] Clean up in the end.")
	sqlDB, _ := dao.GetDbEngine().DB()
	defer sqlDB.Close()

	// Drop Table
	for _, val := range getModels() {
		dao.GetDbEngine().Migrator().DropTable(val)
	}
}

// Run Before a Test
func (t *SuiteTest) SetupTest() {
	//println("***SetupTest, Run Before a Test")
}

// Run After a Test
func (t *SuiteTest) TearDownTest() {
	//println("***TearDownTest, Run After a Test")
}
