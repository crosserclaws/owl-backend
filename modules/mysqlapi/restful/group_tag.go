package restful

import (
	"github.com/fwtpe/owl/common/gin/mvc"
	"github.com/fwtpe/owl/common/model"

	dbOwl "github.com/fwtpe/owl/common/db/owl"
)

func listGroupTags(
	p *struct {
		Name   string        `mvc:"query[name]"`
		Paging *model.Paging `mvc:"pageSize[100] pageOrderBy[name#asc]"`
	},
) (*model.Paging, mvc.OutputBody) {
	return p.Paging,
		mvc.JsonOutputBody(
			dbOwl.ListGroupTags(p.Name, p.Paging),
		)
}

func getGroupTagById(
	p *struct {
		GroupTagId int32 `mvc:"param[group_tag_id]"`
	},
) mvc.OutputBody {
	return mvc.JsonOutputOrNotFound(
		dbOwl.GetGroupTagById(p.GroupTagId),
	)
}
