package rdb

import (
	gormExt "github.com/fwtpe/owl/common/gorm"
	"github.com/fwtpe/owl/modules/mysqlapi/model"
)

func GetAgentConfig(key string) *model.AgentConfigResult {
	var result model.AgentConfigResult
	gormDbExt := gormExt.ToDefaultGormDbExt(
		DbFacade.GormDb.First(&result, "common_config.key = ?", key),
	)

	if gormDbExt.IsRecordNotFound() {
		return nil
	}
	gormDbExt.PanicIfError()

	return &result
}
