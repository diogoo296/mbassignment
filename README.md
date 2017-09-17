# MBAssignment

MessageBird Backend Assignment: a simple API that uses MessageBird API to send SMSes.

## Getting Started

Before anything else, please install or check the current installed version of *Golang* (**>=1.7.4**) and make sure it follows the code organization [described at Golang website](https://golang.org/doc/code.html#Organization).

Now clone this project and install its single dependency (**go-rest-api**):
```
cd $GOPATH/src/github.com
git clone git@bitbucket.org:diogo296/mbassignment.git
cd diogo296/mbassignment
go get github.com/messagebird/go-rest-api
```

Setup your MessageBird API access key in the file *config.json* as the example below, where the *development* key corresponds to the *Mode test* MessageBird key and the *production key* is the *Mode live* key. Notice that it's not required to setup both keys, but the application will not work properly if it's running in a environment with no key set.
```
{
  "MbApiKey": {
    "development": <YOUR_TEST_API_KEY>,
    "production": <YOUR_LIVE_API_KEY>
  }
}
```

At last, build the project by typing the command `go build`, which will generate the bynary file `mbassignment`.

## How to run it

Simply run the binary generated by the build in your terminal with `./mbassignment`. If the project is running properly, you will see a message like this one:
```
2017/09/17 10:56:01 ENV: development
2017/09/17 10:56:01 Server started!
```

The server will run by default in the development environment and thus use MessageBird API test key. To use the live key, please run the API with the *production environment variable*:
```
ENV=producion ./mbassignment
2017/09/17 10:57:10 ENV: production
2017/09/17 10:57:10 Server started!
```

Notice that the application will run in the *development environment* for any `ENV` parameter different from *production*.

## How it works

This API has a single route, `/messages`, which makes a **POST** request to MessageBird API at the URI `https://rest.messagebird.com/messages`. This route expects a json payload with the following parameters:

| Parameter  | Type   | Description |
|------------|--------|-------------|
| recipient  | string | The recipient msisdn. **Required**. |
| originator | string | The sender of the message. This can be a telephone number (including country code) or an alphanumeric string. In case of an alphanumeric string, the maximum length is 11 characters. **Required**. |
| message    | string | The body of the SMS message. **Required**. |

Here is an example of a *curl* request to a running server instance of this API:
```
curl -X POST https://localhost:8080/messages \
-H 'Content-Type: application/json' \
-d '{"recipient":"31612345678", "originator":"YourName", "body": "Test message"}'
```
A successful request should receive as a response a MessageBird message object like the one below:
```
[{"Id":"5e404e4fef8249da98bdd5cfcc0ce2d6","HRef":"https://rest.messagebird.com/messages/5e404e4fef8249da98bdd5cfcc0ce2d6","Direction":"mt","Type":"sms","Originator":"YourName","Body":"Test message","Reference":"","Validity":null,"Gateway":10,"TypeDetails":{},"DataCoding":"plain","MClass":1,"ScheduledDatetime":null,"CreatedDatetime":"2017-09-17T14:42:30Z","Recipients":{"TotalCount":1,"TotalSentCount":1,"TotalDeliveredCount":0,"TotalDeliveryFailedCount":0,"Items":[{"Recipient":31612345678,"Status":"sent","StatusDatetime":"2017-09-17T14:42:30Z"}]},"Errors":null}]
```

## Characters and message length

Considering the concatenated SMS limitations and differences between GSM 03.38 characters and unicode characters, the number of messages sent (and thus the number of messages charged) follows the table described [in this article](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message-).

For more about GSM 03.38 characters and which ones form this set counts as 2 characters in a message, check out [this article](https://support.messagebird.com/hc/en-us/articles/208739765-Which-special-characters-count-as-two-characters-in-a-text-message-).