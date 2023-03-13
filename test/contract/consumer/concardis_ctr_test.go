package main

import (
	"context"
	"errors"
	"fmt"
	auzerolog "github.com/StephanHCB/go-autumn-logging-zerolog"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/concardis"
	"github.com/eurofurence/reg-payment-cncrd-adapter/internal/repository/config"
	"github.com/pact-foundation/pact-go/dsl"
	"log"
	"testing"
)

// contract test consumer side

func TestConcardisApiClient(t *testing.T) {
	auzerolog.SetupPlaintextLogging()

	// Create Pact connecting to local Daemon
	pact := &dsl.Pact{
		Consumer: "reg_payment_cncrd_adapter",
		Provider: "concardis_api",
		Host:     "localhost",
	}
	defer pact.Teardown()

	// set up basic configuration
	config.LoadTestingConfigurationFromPathOrAbort("../../resources/testconfig.yaml")

	// prepare dtos

	// STEP 1: create
	createRequest := concardis.PaymentLinkCreateRequest{
		Title:       "Convention Registration",
		Description: "Please pay for your registration",
		PSP:         1,
		ReferenceId: "220118-150405-000004",
		OrderId:     "220118-150405-000004",
		Purpose:     "EF 2022 REG 000004",
		Amount:      10550,
		VatRate:     19.0,
		Currency:    "EUR",
		SKU:         "REG2022V01AT000004",
		Email:       "test@example.com",
	}
	createRequestSampleBody := `title=Convention+Registration&description=Please+pay+for+your+registration&psp=1&` +
		`referenceId=220118-150405-000004&merchantOrderId=220118-150405-000004&purpose=EF+2022+REG+000004&amount=10550&vatRate=19.0&currency=EUR&` +
		`sku=REG2022V01AT000004&preAuthorization=0&reservation=0&` +
		// `fields[terms][active]=0&fields[terms][mandatory]=0&fields[privacy_policy][active]=0&fields[privacy_policy][mandatory]=0&` +
		`fields[email][mandatory]=1&fields[email][defaultValue]=test@example.com&ApiSignature=omitted`
	createRequestResponse := `{
  "status": "success",
  "data": [
    {
      "id": 42,
      "status": "waiting",
      "hash": "77871ec636d239b136e4d79d32c376ad",
      "referenceId": "220118-150405-000004",
      "link": "http://localhost/some/pay/link",
      "invoices": [],
      "merchantOrderId": "220118-150405-000004",
      "preAuthorization": false,
      "name": "",
      "title": "Convention Registration",
      "description": "Please pay for your registration",
      "buttonText": "",
      "api": true,
      "fields": {
        "header": {
          "active": true,
          "mandatory": false,
          "names": {
            "de": "Kontaktdaten"
          }
        }
      },
      "psp": [],
      "pm": [],
      "purpose": {
        "1": "EF 2022 REG 000004"
      },
      "amount": 10550,
      "currency": "EUR",
      "vatRate": "19.0",
      "sku": "REG2022V01AT000004",
      "subscriptionState": false,
      "subscriptionInterval": "",
      "subscriptionPeriod": "",
      "subscriptionPeriodMinAmount": "",
      "subscriptionCancellationInterval": "",
      "createdAt": 1665838673,
      "requestId": 7489378
    }
  ]
}`

	// STEP 2: query after it has been used
	queryRequestSampleBody := `ApiSignature=omitted`
	queryRequestResponse := `{
  "status": "success",
  "data": [
    {
      "id": 424242,
      "status": "confirmed",
      "hash": "77871ec636d239b136e4d79d32c376ad",
      "referenceId": "220118-150405-000004",
      "link": "http://localhost/some/pay/link",
      "invoices": [
        {
          "number": "EF 2022 REG 000004",
          "products": [
            {
              "name": "EF 2022 REG 000004",
              "price": 10550,
              "quantity": 1,
              "sku": "REG2022V01AT000004",
              "vatRate": null
            }
          ],
          "amount": 10550,
          "discount": {
            "code": null,
            "amount": 0,
            "percentage": null
          },
          "shippingAmount": null,
          "currency": "EUR",
          "test": 1,
          "referenceId": "220118-150405-000004",
          "paymentRequestId": 42,
          "paymentLink": {
            "hash": "77871ec636d239b136e4d79d32c376ad",
            "referenceId": "220118-150405-000004",
            "email": ""
          },
          "transactions": [
            {
              "id": 777777,
              "uuid": "b9bee580",
              "amount": 10550,
              "referenceId": "220118-150405-000004",
              "time": "2022-10-15 15:50:20",
              "status": "confirmed",
              "lang": "de",
              "psp": "ConCardis_PayEngine_3",
              "pspId": 29,
              "mode": "TEST",
              "metadata": [],
              "contact": {
                "id": 888888,
                "uuid": "50659dad",
                "title": "",
                "firstname": "",
                "lastname": "",
                "company": "",
                "street": "",
                "zip": "",
                "place": "",
                "country": "",
                "countryISO": "",
                "phone": "",
                "email": "",
                "date_of_birth": "",
                "delivery_title": "",
                "delivery_firstname": "",
                "delivery_lastname": "",
                "delivery_company": "",
                "delivery_street": "",
                "delivery_zip": "",
                "delivery_place": "",
                "delivery_country": "",
                "delivery_countryISO": ""
              },
              "subscription": null,
              "pageUuid": null,
              "payment": {
                "brand": "visa"
              }
            }
          ],
          "custom_fields": [
            {
              "type": "header",
              "name": "Kontaktdaten",
              "value": "Kontaktdaten"
            }
          ]
        }
      ],
      "merchantOrderId": "220118-150405-000004",
      "preAuthorization": false,
      "name": "",
      "title": "Convention Registration",
      "description": "Please pay for your registration",
      "buttonText": "",
      "api": true,
      "fields": {
        "header": {
          "active": true,
          "mandatory": false,
          "names": {
            "de": "Kontaktdaten"
          }
        }
      },
      "psp": [],
      "pm": [],
      "purpose": {
        "1": "EF 2022 REG 000004"
      },
      "amount": 10550,
      "currency": "EUR",
      "vatRate": 19.0,
      "sku": "REG2022V01AT000004",
      "subscriptionState": false,
      "subscriptionInterval": "",
      "subscriptionPeriod": "",
      "subscriptionPeriodMinAmount": 0,
      "subscriptionCancellationInterval": "",
      "createdAt": 1665838673
    }
  ]
}`
	// STEP 3: delete
	deleteRequestSampleBody := `ApiSignature=omitted`

	// Pass in test case (consumer side)
	// This uses the repository on the consumer side to make the http call, should be as low level as possible
	var test = func() (err error) {
		ctx := auzerolog.AddLoggerToCtx(context.Background())

		// override configuration with pact server url
		config.Configuration().Service.ConcardisDownstream = fmt.Sprintf("http://localhost:%d", pact.Server.Port)

		// set up downstream client
		err = concardis.Create()
		if err != nil {
			return err
		}
		client := concardis.Get()
		// pact does not support regex matchers for x-www-form-urlencoded, when it eventually does, we could use the actual signature
		concardis.FixedSignatureValue = "omitted"

		// STEP 1: create a new payment link
		created, err := client.CreatePaymentLink(ctx, createRequest)
		if err != nil {
			return err
		}
		if created.ID != 42 || created.ReferenceID != "220118-150405-000004" || created.Link != "http://localhost/some/pay/link" {
			return errors.New("unexpected create response")
		}

		// STEP 2: read the payment link again after use
		read, err := client.QueryPaymentLink(ctx, created.ID)
		if err != nil {
			return err
		}
		if read.Invoices[0].PaymentRequestId != 42 || read.ReferenceID != "220118-150405-000004" || read.Status != "confirmed" {
			return errors.New("unexpected query response")
		}

		// STEP 3: delete the payment link (wouldn't normally work after use)
		err = client.DeletePaymentLink(ctx, created.ID)
		if err != nil {
			return err
		}

		return nil
	}

	// Set up our expected interactions.
	pact.
		AddInteraction().
		UponReceiving("A request to create a payment link").
		// pact does not support regex matchers for x-www-form-urlencoded
		WithRequest(dsl.Request{
			Method:  "POST",
			Path:    dsl.String("/v1.0/Invoice/"),
			Query:   dsl.MapMatcher{"instance": dsl.String("myinstance")},
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/x-www-form-urlencoded")},
			Body:    dsl.String(createRequestSampleBody),
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.String(createRequestResponse),
		})

	pact.
		AddInteraction().
		UponReceiving("A request to read the same payment link").
		WithRequest(dsl.Request{
			Method:  "GET",
			Path:    dsl.String("/v1.0/Invoice/42/"),
			Query:   dsl.MapMatcher{"instance": dsl.String("myinstance")},
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/x-www-form-urlencoded")},
			Body:    dsl.String(queryRequestSampleBody),
		}).
		WillRespondWith(dsl.Response{
			Status:  200,
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/json")},
			Body:    dsl.String(queryRequestResponse),
		})

	pact.
		AddInteraction().
		UponReceiving("A request to delete the same payment link").
		WithRequest(dsl.Request{
			Method:  "DELETE",
			Path:    dsl.String("/v1.0/Invoice/42/"),
			Query:   dsl.MapMatcher{"instance": dsl.String("myinstance")},
			Headers: dsl.MapMatcher{"Content-Type": dsl.String("application/x-www-form-urlencoded")},
			Body:    dsl.String(deleteRequestSampleBody),
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
		})

	// Run the test, verify it did what we expected and capture the contract (writes a test log to logs/pact.log)
	if err := pact.Verify(test); err != nil {
		log.Fatalf("Error on Verify: %v", err)
	}

	// now write out the contract json (by default it goes to subdirectory pacts)
	if err := pact.WritePact(); err != nil {
		log.Fatalf("Error on pact write: %v", err)
	}
}
