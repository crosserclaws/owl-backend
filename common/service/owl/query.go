package owl

import (
	"net/http"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	gt "gopkg.in/h2non/gentleman.v2"

	"github.com/fwtpe/owl/common/db"
	oHttp "github.com/fwtpe/owl/common/http"
	"github.com/fwtpe/owl/common/http/client"
	"github.com/fwtpe/owl/common/json"
	ojson "github.com/fwtpe/owl/common/json"
	model "github.com/fwtpe/owl/common/model/owl"
)

type QueryServiceConfig struct {
	*oHttp.RestfulClientConfig
}

type QueryService interface {
	LoadQueryByUuid(uuid.UUID) *model.Query
	CreateOrLoadQuery(*model.Query)
	VacuumQueryObjects(int) *ResultOfVacuumQueryObjects
}

func NewQueryService(config QueryServiceConfig) QueryService {
	newClient := oHttp.NewApiService(config.RestfulClientConfig).NewClient()

	return &queryServiceImpl{
		loadQueryByUuid:    newClient.Get().AddPath("/api/v1/owl/query-object"),
		createOrLoadQuery:  newClient.Post().AddPath("/api/v1/owl/query-object"),
		vacuumQueryObjects: newClient.Post().AddPath("/api/v1/owl/query-object/vacuum"),
	}
}

type queryServiceImpl struct {
	loadQueryByUuid    *gt.Request
	createOrLoadQuery  *gt.Request
	vacuumQueryObjects *gt.Request
}

func (s *queryServiceImpl) LoadQueryByUuid(uuid uuid.UUID) *model.Query {
	req := s.loadQueryByUuid.Clone().
		AddPath("/" + uuid.String())

	resp := client.ToGentlemanReq(req).SendAndMustMatch(
		func(resp *gt.Response) error {
			switch resp.StatusCode {
			case http.StatusOK, http.StatusNotFound:
				return nil
			}

			return errors.Errorf(client.ToGentlemanResp(resp).ToDetailString())
		},
	)

	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	jsonQuery := &struct {
		Uuid    ojson.Uuid `json:"uuid"`
		NamedId string     `json:"feature_name"`

		Content    ojson.VarBytes `json:"content"`
		Md5Content ojson.Bytes16  `json:"md5_content"`

		CreationTime ojson.JsonTime `json:"creation_time"`
		AccessTime   ojson.JsonTime `json:"access_time"`
	}{}

	client.ToGentlemanResp(resp).MustBindJson(jsonQuery)

	return &model.Query{
		Uuid:    db.DbUuid(jsonQuery.Uuid),
		NamedId: jsonQuery.NamedId,

		Content:    []byte(jsonQuery.Content),
		Md5Content: db.Bytes16(jsonQuery.Md5Content),

		CreationTime: time.Time(jsonQuery.CreationTime),
		AccessTime:   time.Time(jsonQuery.AccessTime),
	}
}

// Loads object of query or creating one.
//
// Any error would be expressed by panic.
func (s *queryServiceImpl) CreateOrLoadQuery(query *model.Query) {
	req := s.createOrLoadQuery.Clone().
		JSON(map[string]interface{}{
			"feature_name": query.NamedId,
			"content":      ojson.VarBytes(query.Content),
			"md5_content":  ojson.Bytes16(query.Md5Content),
		})

	resp := client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK)

	jsonBody := &struct {
		Uuid         ojson.Uuid     `json:"uuid"`
		CreationTime ojson.JsonTime `json:"creation_time"`
		AccessTime   ojson.JsonTime `json:"access_time"`
	}{}

	client.ToGentlemanResp(resp).MustBindJson(jsonBody)

	query.Uuid = db.DbUuid(jsonBody.Uuid)
	query.CreationTime = time.Time(jsonBody.CreationTime)
	query.AccessTime = time.Time(jsonBody.AccessTime)
}

type ResultOfVacuumQueryObjects struct {
	BeforeTime   json.JsonTime `json:"before_time"`
	AffectedRows int           `json:"affected_rows"`
}

func (r *ResultOfVacuumQueryObjects) GetBeforeTime() time.Time {
	return time.Time(r.BeforeTime)
}

// Vacuums out-dated query objects(by access time)
//
// Any error would be expressed by panic.
func (s *queryServiceImpl) VacuumQueryObjects(forDays int) *ResultOfVacuumQueryObjects {
	req := s.vacuumQueryObjects.Clone().
		SetQueryParams(map[string]string{
			"for_days": strconv.Itoa(forDays),
		})

	result := &ResultOfVacuumQueryObjects{}

	client.ToGentlemanResp(
		client.ToGentlemanReq(req).SendAndStatusMustMatch(http.StatusOK),
	).MustBindJson(result)

	return result
}
