package rdb

import (
	"fmt"

	commonDb "github.com/fwtpe/owl-backend/common/db"
	f "github.com/fwtpe/owl-backend/common/db/facade"
	commonNqmDb "github.com/fwtpe/owl-backend/common/db/nqm"
	commonOwlDb "github.com/fwtpe/owl-backend/common/db/owl"
	log "github.com/fwtpe/owl-backend/common/logruslog"
	apiModel "github.com/fwtpe/owl-backend/common/model/mysqlapi"

	bossdb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/boss"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/cmdb"
	graphdb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/graph"
	"github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/hbsdb"
	apiOwlDb "github.com/fwtpe/owl-backend/modules/mysqlapi/rdb/owl"
)

const (
	DB_PORTAL = "portal"
	DB_GRAPH  = "graph"
	DB_BOSS   = "boss"
)

type DbHolder struct {
	facades map[string]*f.DbFacade
}

func (self *DbHolder) setDb(dbname string, facade *f.DbFacade) {
	self.facades[dbname] = facade
}
func (self *DbHolder) releaseDb(dbname string) {
	if facade, ok := self.facades[dbname]; ok {
		facade.Release()
		delete(self.facades, dbname)
	}
}
func (self *DbHolder) Diagnose(dbname string) *apiModel.Rdb {
	facade, ok := self.facades[dbname]

	if !ok {
		return nil
	}

	return DiagnoseRdb(facade.GetDbConfig().Dsn, facade.SqlDb)
}

var GlobalDbHolder *DbHolder = &DbHolder{
	facades: make(map[string]*f.DbFacade),
}

var logger = log.NewDefaultLogger("INFO")

var DbFacade *f.DbFacade

func InitPortalRdb(dbConfig *commonDb.DbConfig) {
	openRdb(
		DB_PORTAL, dbConfig,
		func(facade *f.DbFacade) {
			facade.SetReleaseCallback(func() {
				commonNqmDb.DbFacade = nil
				commonOwlDb.DbFacade = nil
				apiOwlDb.DbFacade = nil

				hbsdb.DbFacade = nil
				hbsdb.DB = nil
			})

			DbFacade = facade

			/**
			 * Protal database
			 */
			commonNqmDb.DbFacade = DbFacade
			commonOwlDb.DbFacade = DbFacade
			apiOwlDb.DbFacade = DbFacade
			cmdb.DbFacade = DbFacade

			hbsdb.DbFacade = DbFacade
			hbsdb.DB = DbFacade.SqlDb
			// :~)
		},
	)
}
func InitGraphRdb(dbConfig *commonDb.DbConfig) {
	openRdb(
		DB_GRAPH, dbConfig,
		func(facade *f.DbFacade) {
			facade.SetReleaseCallback(func() {
				graphdb.DbFacade = nil
			})

			graphdb.DbFacade = facade
		},
	)
}
func InitBossRdb(dbConfig *commonDb.DbConfig) {
	openRdb(
		DB_BOSS, dbConfig,
		func(facade *f.DbFacade) {
			facade.SetReleaseCallback(func() {
				bossdb.DbFacade = nil
			})

			bossdb.DbFacade = facade
		},
	)
}

type displayDbConfig commonDb.DbConfig

func (c *displayDbConfig) String() string {
	return fmt.Sprintf("DSN: [%s]. Max Idle: [%d]", hidePasswordOfDsn(c.Dsn), c.MaxIdle)
}

func openRdb(dbName string, dbConfig *commonDb.DbConfig, facadeCallback func(*f.DbFacade)) {
	newFacade := &f.DbFacade{}
	GlobalDbHolder.setDb(dbName, newFacade)

	logger.Infof("Open RDB: %s ...", (*displayDbConfig)(dbConfig))

	err := newFacade.Open(dbConfig)
	if err != nil {
		logger.Warnf("Open database error: %v", err)
	}

	facadeCallback(newFacade)

	logger.Info("[FINISH] Open RDB.")
}

func ReleaseAllRdb() {
	logger.Info("Release RDB resources...")

	GlobalDbHolder.releaseDb(DB_PORTAL)
	GlobalDbHolder.releaseDb(DB_GRAPH)
	GlobalDbHolder.releaseDb(DB_BOSS)

	logger.Info("[FINISH] Release RDB resources.")
}
