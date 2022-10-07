package acceptance

import (
	"encoding/json"
	"fmt"
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

func tstPerformGet(relativeUrlWithLeadingSlash string, bearerToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPut(relativeUrlWithLeadingSlash string, requestBody string, bearerToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPut, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	request.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return tstWebResponseFromResponse(response)
}

func tstPerformPost(relativeUrlWithLeadingSlash string, requestBody string, bearerToken string) tstWebResponse {
	request, err := http.NewRequest(http.MethodPost, ts.URL+relativeUrlWithLeadingSlash, strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(err)
	}
	if bearerToken != "" {
		request.Header.Set(headers.Authorization, "Bearer "+bearerToken)
	}
	request.Header.Set(headers.ContentType, media.ContentTypeApplicationJson)
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

func tstBuildValidPaymentLinkRequest(testcase string) cncrdapi.PaymentLink {
	return cncrdapi.PaymentLink{
		Title:       "some title",
		Description: "some description",
		ReferenceId: fmt.Sprintf("144823ad-%s", testcase),
		Purpose:     "some purpose",
		AmountDue:   390,
		AmountPaid:  0,
		Currency:    "EUR",
		VatRate:     19.0,
		Link:        "",
	}
}
