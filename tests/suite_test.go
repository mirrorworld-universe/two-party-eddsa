package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"main/global"
	"main/middleware/dao"
	validator2 "main/middleware/validator"
	"main/model/db"
	"main/routes"
	"net/http"
	"os"
	"testing"
)

type SuiteTest struct {
	suite.Suite
	router *gin.Engine
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

	// extra config
	gin.SetMode(gin.TestMode)

	// custom validators
	validator2.SetupValidators()

	// setup router
	router := routes.NewRouter()

	// create service
	srv := &http.Server{
		Addr:    global.Config.Base.Port,
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			global.Logger.Error("Gin server start error:", err.Error())
			panic(err.Error())
		}
	}()

	t.router = router
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
