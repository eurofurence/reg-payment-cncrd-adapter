package acceptance

import (
	"encoding/json"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/api/v1/cncrdapi"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/web/util/media"
	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// placing these here because they are package global

type tstWebResponse struct {
	status      int
	body        string
	contentType string
	location    string
}

func tstWebResponseFromResponse(response *http.Response) tstWebResponse {
	status := response.StatusCode
	ct := ""
	if val, ok := response.Header[headers.ContentType]; ok {
		ct = val[0]
	}
	loc := ""
	if val, ok := response.Header[headers.Location]; ok {
		loc = val[0]
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponse{
		status:      status,
		body:        string(body),
		contentType: ct,
		location:    loc,
	}
}

func tstPerformGet(relativeUrlWithLeadingSlash string, apiToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	if apiToken != "" {
		request.Header.Set(media.HeaderXApiKey, apiToken)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPost(relativeUrlWithLeadingSlash string, requestBody string, apiToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPost, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	if apiToken != "" {
		request.Header.Set(media.HeaderXApiKey, apiToken)
	}
	request.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformDelete(relativeUrlWithLeadingSlash string, apiToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodDelete, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	if apiToken != "" {
		request.Header.Set(media.HeaderXApiKey, apiToken)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstRenderJson(v interface{}) string {
	representationBytes, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}
	return string(representationBytes)
}

// tip: dto := &XyzDto{}
func tstParseJson(body string, dto interface{}) {
	err := json.Unmarshal([]byte(body), dto)
	if err != nil {
		log.Fatal(err)
	}
}

func tstRequireErrorResponse(t *testing.T, response tstWebResponse, expectedStatus int, expectedMessage string, expectedDetails interface{}) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	errorDto := cncrdapi.ErrorDto{}
	tstParseJson(response.body, &errorDto)
	require.Equal(t, expectedMessage, errorDto.Message, "unexpected error code")
	expectedDetailsStr, ok := expectedDetails.(string)
	if ok && expectedDetailsStr != "" {
		require.EqualValues(t, url.Values{"details": []string{expectedDetailsStr}}, errorDto.Details, "unexpected error details")
	}
	expectedDetailsUrlValues, ok := expectedDetails.(url.Values)
	if ok {
		require.EqualValues(t, expectedDetailsUrlValues, errorDto.Details, "unexpected error details")
	}
}

func tstRequirePaymentLinkResponse(t *testing.T, response tstWebResponse, expectedStatus int, expectedBody cncrdapi.PaymentLinkDto) {
	require.Equal(t, expectedStatus, response.status, "unexpected http response status")
	actualBody := cncrdapi.PaymentLinkDto{}
	tstParseJson(response.body, &actualBody)
	require.EqualValues(t, expectedBody, actualBody)
}

func tstRequireConcardisRecording(t *testing.T, expectedEntries ...string) {
	actual := concardisMock.Recording()
	require.Equal(t, len(expectedEntries), len(actual))
	for i := range expectedEntries {
		require.Equal(t, expectedEntries[i], actual[i])
	}
}

// --- data ---

func tstBuildValidPaymentLinkRequest() cncrdapi.PaymentLinkRequestDto {
	return cncrdapi.PaymentLinkRequestDto{
		DebitorId: 1,
		AmountDue: 390,
		Currency:  "EUR",
		VatRate:   19.0,
	}
}

func tstBuildValidPaymentLink() cncrdapi.PaymentLinkDto {
	return cncrdapi.PaymentLinkDto{
		Title:       "some page title",
		Description: "some page description",
		ReferenceId: "221216-122218-000001",
		Purpose:     "some payment purpose",
		AmountDue:   390,
		AmountPaid:  0,
		Currency:    "EUR",
		VatRate:     19.0,
		Link:        "http://localhost:1111/some/paylink",
	}
}

func tstBuildValidPaymentLinkGetResponse() cncrdapi.PaymentLinkDto {
	return cncrdapi.PaymentLinkDto{
		ReferenceId: "221216-122218-000001",
		Purpose:     "some payment purpose",
		AmountDue:   390,
		AmountPaid:  0,
		Currency:    "EUR",
		Link:        "http://localhost:1111/some/paylink",
	}
}
