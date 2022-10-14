package tests

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"main/model/db"
	"os"
	"testing"
)

type SuiteTest struct {
	suite.Suite
	db *gorm.DB
}

func TestSuite(t *testing.T) {
	os.Setenv("ENV", "tests")
	defer os.Unsetenv("ENV")

	suite.Run(t, new(SuiteTest))
}

func getModels() []interface{} {
	return []interface{}{
		&db.MPCWallet{},
	}
}

// Setup db value
func (t *SuiteTest) SetupSuite() {
	//main.InitConfig()
	//main.InitDB()
	//main.InitLogger()

	SetupAll()

}

// Run After All Test Done
func (t *SuiteTest) TearDownSuite() {
	sqlDB, _ := t.db.DB()
	defer sqlDB.Close()

	// Drop Table
	//for _, val := range getModels() {
	//	t.db.Migrator().DropTable(val)
	//}
}

// Run Before a Test
func (t *SuiteTest) SetupTest() {
	println("SetupTest, Run Before a Test")
}

// Run After a Test
func (t *SuiteTest) TearDownTest() {
	println("TearDownTest, Run After a Test")
}
