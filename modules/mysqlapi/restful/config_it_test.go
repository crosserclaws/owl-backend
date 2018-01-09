package restful

import (
	"net/http"

	json "github.com/fwtpe/owl/common/json"
	ogko "github.com/fwtpe/owl/common/testing/ginkgo"
	tHttp "github.com/fwtpe/owl/common/testing/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("[Intg] Test getAgentConfig", itSkipOnPortal.PrependBeforeEach(func() {
	BeforeEach(func() {
		inPortalTx(
			"INSERT INTO common_config(`key`, `value`)" +
				"VALUES('getAgentConfig', 'https://example.com/Cepave/getAgentConfig.git')",
		)
	})

	AfterEach(func() {
		inPortalTx(
			"DELETE FROM common_config WHERE `key` = 'getAgentConfig'",
		)
	})

	DescribeTable("when value",
		func(key string, expectedJson string) {
			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewClient().Get("api/v1/agent/config?key=" + key),
			)
			Expect(result).To(ogko.MatchHttpStatus(http.StatusOK))

			jsonBody := result.GetBodyAsJson()
			jsonResult := json.MarshalPrettyJSON(jsonBody)
			Expect(jsonResult).To(MatchJSON(expectedJson))
		},
		Entry("can be found by key", "getAgentConfig",
			`{
				"key": "getAgentConfig",
				"value": "https://example.com/Cepave/getAgentConfig.git"
			}`,
		),
	)

	DescribeTable("404 status when",
		func(key string) {
			result := tHttp.NewResponseResultBySling(
				httpClientConfig.NewClient().Get("api/v1/agent/config?key=" + key),
			)
			Expect(result).To(ogko.MatchHttpStatus(http.StatusNotFound))
		},
		Entry("empty key parameter", ""),
		Entry("value is not found", "Non-existent-key"),
	)
}))
