
# create a payment link example

```
apiSignature=`echo -n "title=Test&description=Testdescription&psp=1&referenceId=test&purpose=This+is+a+test&amount=200&vatRate=7.7&currency=CHF&sku=P01122000&preAuthorization=0&reservation=0" | openssl dgst -sha256 -hmac "INSTANCE_API_SECRET" -binary | openssl enc -base64`
curl --request POST "https://api.pay-link.eu/v1.0/Invoice/?instance=INSTANCE_NAME" \
  --data-urlencode "title=Test" \
  --data-urlencode "description=Testdescription" \
  --data-urlencode "psp=1" \
  --data-urlencode "referenceId=test" \
  --data-urlencode "purpose=This is a test" \
  --data-urlencode "amount=200" \
  --data-urlencode "vatRate=7.7" \
  --data-urlencode "currency=CHF" \
  --data-urlencode "sku=P01122000" \
  --data-urlencode "preAuthorization=0" \
  --data-urlencode "reservation=0" \
  --data-urlencode "ApiSignature=$apiSignature"
```

response:

```
{
    "status": "success",
    "data": [
        {
            "id": 1,
            "hash": "382c85eab7a86278e3c3b06a23af2358",
            "referenceId": "Order number of my online shop application",
            "link": "https://demo.pay-link.eu/?payment=382c85eab7a86278e3c3b06a23af2358",
            "invoices": [],
            "preAuthorization": 0,
            "reservation": 0,
            "name": "Online-Shop payment #001",
            "api": true,
            "fields": {
                "title": {
                    "active": true,
                    "mandatory": true
                },
                "forename": {
                    "active": true,
                    "mandatory": true
                },
                "surname": {
                    "active": true,
                    "mandatory": true
                },
                "company": {
                    "active": true,
                    "mandatory": true
                },
                "street": {
                    "active": false,
                    "mandatory": false
                },
                "postcode": {
                    "active": false,
                    "mandatory": false
                },
                "place": {
                    "active": false,
                    "mandatory": false
                },
                "country": {
                    "active": true,
                    "mandatory": true
                },
                "phone": {
                    "active": false,
                    "mandatory": false
                },
                "email": {
                    "active": true,
                    "mandatory": true
                },
                "date_of_birth": {
                    "active": false,
                    "mandatory": false
                },
                "terms": {
                    "active": true,
                    "mandatory": true
                },
                "privacy_policy": {
                    "active": true,
                    "mandatory": true
                },
                "custom_field_1": {
                    "active": true,
                    "mandatory": true,
                    "names": {
                        "de": "This is a field",
                        "en": "This is a field",
                        "fr": "This is a field",
                        "it": "This is a field"
                    }
                },
                "custom_field_2": {
                    "active": false,
                    "mandatory": false,
                    "names": {
                        "de": "",
                        "en": "",
                        "fr": "",
                        "it": ""
                    }
                },
                "custom_field_3": {
                    "active": false,
                    "mandatory": false,
                    "names": {
                        "de": "",
                        "en": "",
                        "fr": "",
                        "it": ""
                    }
                }
            },
            "psp": 1,
            "pm": [],
            "purpose": "Shop Order #001",
            "amount": 590,
            "vatRate" : 7.7,
            "currency": "CHF",
            "sku": "P01122000",
            "subscriptionState": false,
            "subscriptionInterval": "",
            "subscriptionPeriod": "",
            "subscriptionPeriodMinAmount": "",
            "subscriptionCancellationInterval": "",
            "createdAt": 1418392958
        }
    ]
}
```

# read a payment link example

```
apiSignature=`echo -n "" | openssl dgst -sha256 -hmac "INSTANCE_API_SECRET" -binary | openssl enc -base64`
curl --request GET "https://api.pay-link.eu/v1.0/Invoice/1/?instance=INSTANCE_NAME" --data-urlencode "ApiSignature=$apiSignature"
```

response:

```
{
    "status": "success",
    "data": [
        {
            "id": 1,
            "status": "(waiting|confirmed|authorized|reserved)",
            "hash": "382c85eab7a86278e3c3b06a23af2358",
            "referenceId": "Order number of my online shop application",
            "link": "https://demo.pay-link.eu/?payment=382c85eab7a86278e3c3b06a23af2358",
            "invoices": [],
            "name": "Online-Shop payment #001",
            "api": true,
            "fields": {
                "title": {
                    "active": true,
                    "mandatory": true
                },
                "forename": {
                    "active": true,
                    "mandatory": true
                },
                "surname": {
                    "active": true,
                    "mandatory": true
                },
                "company": {
                    "active": true,
                    "mandatory": true
                },
                "street": {
                    "active": false,
                    "mandatory": false
                },
                "postcode": {
                    "active": false,
                    "mandatory": false
                },
                "place": {
                    "active": false,
                    "mandatory": false
                },
                "country": {
                    "active": true,
                    "mandatory": true
                },
                "phone": {
                    "active": false,
                    "mandatory": false
                },
                "email": {
                    "active": true,
                    "mandatory": true
                },
                "date_of_birth": {
                    "active": false,
                    "mandatory": false
                },
                "terms": {
                    "active": true,
                    "mandatory": true
                },
                "custom_field_1": {
                    "active": true,
                    "mandatory": true,
                    "names": {
                        "de": "This is a field",
                        "en": "This is a field",
                        "fr": "This is a field",
                        "it": "This is a field"
                    }
                },
                "custom_field_2": {
                    "active": false,
                    "mandatory": false,
                    "names": {
                        "de": "",
                        "en": "",
                        "fr": "",
                        "it": ""
                    }
                },
                "custom_field_3": {
                    "active": false,
                    "mandatory": false,
                    "names": {
                        "de": "",
                        "en": "",
                        "fr": "",
                        "it": ""
                    }
                }
            },
            "psp": 1,
            "purpose": "Shop Order #001",
            "amount": 590,
            "currency": "CHF",
            "subscriptionState": false,
            "subscriptionInterval": "",
            "subscriptionPeriod": "",
            "subscriptionPeriodMinAmount": 0,
            "subscriptionCancellationInterval": "",
            "createdAt": 1418392958
        }
    ]
}
```

# delete payment link

```
apiSignature=`echo -n "" | openssl dgst -sha256 -hmac "INSTANCE_API_SECRET" -binary | openssl enc -base64`
curl --request DELETE "https://api.pay-link.eu/v1.0/Invoice/1/?instance=INSTANCE_NAME" --data-urlencode "ApiSignature=$apiSignature"
```

response:

```
{
    "status": "success",
    "data": [
        {
            "id": 1
        }
    ]
}
```


# read transactions

```
apiSignature=`echo -n "filterDatetimeUtcGreaterThan=2020-01-20+14%3A55%3A00&filterDatetimeUtcLessThan=2020-01-25+19%3A00%3A50&offset=30&limit=2" | openssl dgst -sha256 -hmac "YOUR_SECRET" -binary | openssl enc -base64`
curl --request GET "https://api.pay-link.eu/v1.0/Transaction//?instance=YOUR_INSTANCE_NAME" \
--data-urlencode "filterDatetimeUtcGreaterThan=2020-01-20 14:55:00" \
--data-urlencode "filterDatetimeUtcLessThan=2020-01-25 19:00:50" \
--data-urlencode "offset=30" \
--data-urlencode "limit=2" \
--data-urlencode "ApiSignature=$apiSignature"
```

response:

```
{ 
    "status":"success",
    "data":[ 
        { 
            "id":1,
            "uuid":"f384000b",
            "status":"authorized",
            "time":"2020-01-20 15:56:02",
            "lang":"de",
            "pageUuid":"892dcf5c",
            "payment":{ 
                "brand":"visa"
            },
            "psp":"ConCardis_PayEngine",
            "pspId":9,
            "mode":"LIVE",
            "referenceId":"ORDER#321",
            "invoice":{ 
                "currencyAlpha3":"CHF",
                "products":[ 
                    { 
                        "quantity":1,
                        "name":"Hoodie",
                        "amount":5900
                    },
                    { 
                        "quantity":2,
                        "name":"T-Shirt",
                        "amount":2850
                    }
                ],
                "discount":{ 
                    "code":"SUPERSALE",
                    "percentage":33,
                    "amount":3828
                },
                "shippingAmount":150000,
                "totalAmount":9272,
                "customFields":{ 
                    "20":{ 
                        "name":"Additional Information",
                        "value":"Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."
                    }
                }
            },
            "contact":{ 
                "id":16,
                "uuid":"9c9c0282",
                "company":"FeelGood",
                "title":"mister",
                "firstname":"Hans",
                "lastname":"Schmid",
                "street":"Dorstrasse 11",
                "zip":"1234",
                "place":"Habern",
                "country":"Schweiz",
                "countryISO":"CH",
                "date_of_birth":"10.03.1989",
                "email":"info@example.com",
                "phone":"+413312345678",
                "delivery_company":"NewFashion",
                "delivery_title":"miss",
                "delivery_firstname":"Fabienne",
                "delivery_lastname":"Muster",
                "delivery_street":"Seestrasse 39",
                "delivery_zip":"4321",
                "delivery_place":"Nechters",
                "delivery_country":"Schweiz",
                "delivery_countryISO":"CH"
            }
        },
        { 
            "id":223,
            "uuid":"d2571806",
            "status":"confirmed",
            "time":"2020-01-20 16:25:58",
            "lang":"de",
            "pageUuid":"892dcf5c",
            "payment":{ 
                "brand":"visa"
            },
            "psp":"ConCardis_PayEngine",
            "pspId":9,
            "mode":"LIVE",
            "referenceId":"ORDER#322",
            "invoice":{ 
                "currencyAlpha3":"CHF",
                "products":[ 
                    { 
                        "name":"T-Shirt",
                        "quantity":1,
                        "amount":2850
                    }
                ],
                "discount":null,
                "shippingAmount":150000,
                "totalAmount":4350,
                "customFields":null
            },
            "contact":{ 
                "id":16,
                "uuid":"9c9c0282",
                "company":"",
                "title":"mister",
                "lastname":"Muster",
                "firstname":"Hans",
                "street":"Dorfstrasse 49",
                "zip":"1234",
                "place":"Habern",
                "country":"Schweiz",
                "countryISO":"CH",
                "date_of_birth":"10.03.1989",
                "email":"info@example.com",
                "phone":"",
                "delivery_company":"",
                "delivery_title":"0",
                "delivery_firstname":"",
                "delivery_lastname":"",
                "delivery_street":"",
                "delivery_zip":"",
                "delivery_place":"",
                "delivery_country":"",
                "delivery_countryISO":""
            }
        }
    ]
}
```

# webhook incoming

```
array(
  'id' => 1,
  'uuid' => '82m09f9',
  'time' => '2014-11-18 13:44:53',
  'status' => 'waiting',
  'lang' => 'en',
  'psp' => 'Test',
  'payment' => array(
    'brand' => 'VISA'
  ),
  'metadata' => Metadata,
  'subscription' => Subscription,
  'invoice' => Invoice,
  'contact' => Contact,
)
```