# MBAssignment

MessageBird Backend Assignment: a simple API that uses MessageBird API to send SMSes.

## 1. Getting Started

Before anything else, please install or check the current installed version of *Golang* (**>=1.7.4**) and make sure it follows the code organization [described at Golang website](https://golang.org/doc/code.html#Organization).

### 1.1. Clone this repository

Clone this project and install its single dependency ([go-rest-api](https://github.com/messagebird/go-rest-api)):
```
cd $GOPATH/src/github.com
git clone git@bitbucket.org:diogo296/mbassignment.git
go get github.com/messagebird/go-rest-api
```

### 1.2. Setup the config file

Go to the project folder and fill the configurations in *config.json* file. This file should looks like the example below:
```
{
  "MbApiKey": {
    "development": <YOUR_TEST_API_KEY>,
    "production": <YOUR_LIVE_API_KEY>
  },
  "MaxCsmsParts": 9
}
```
Here is a quick description about the parameters:
* **MbApiKey:** Your MessageBird REST API key. The *development* entry corresponds to the *Mode test* key and the *production key* is to the *Mode live* key. Notice that it's not required to setup both of them, but the application will not work properly if it's running in a environment with the correspondent key not set.
* **MaxCsmsParts:** Maximum number of concatenated SMS parts. Following MessageBird Dashboard Quick Send interface, it was initially set to 9 but it is not a mandatory parameter and, if removed from the configuration, it **will be set to 255** due to the limitation of number of parts that can be encoded in a CSMS UDH. Must be **bewteen 0 and 255**.

At last, you can check if your configuration file is fine by running the command `go test -run TestLoadConfig`.

### 1.3. Build and Run

Build the project by typing the command `go build`, which will generate the bynary file `mbassignment`.

To run the API, somply type `./mbassignment` in your terminal. If the project is running properly, the API will be up and running at `localhost:8080` and you will see a message like this one:
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

## 2. How it works

This API has a single endpoint, `/messages`, which makes a **POST** request to MessageBird API at the URI `https://rest.messagebird.com/messages`. This route expects a json payload with the following parameters:

| Parameter  | Type   | Description |
|------------|--------|-------------|
| recipient  | string | The recipient msisdn. **Required**. |
| originator | string | The sender of the message. This can be a telephone number (including country code) or an alphanumeric string. In case of an alphanumeric string, the maximum length is 11 characters. **Required**. |
| message    | string | The body of the SMS message. Can be composed by either GSM 03.38 characters or unicode characters and must contain at least one non-whitespace character. **Required**. |

Here is an example of a *curl* request to a running server instance of this API:
```
curl -X POST http://localhost:8080/messages \
-H 'Content-Type: application/json' \
-d '{"recipient":"31612345678", "originator":"YourName", "message": "Test message"}'
```
A successful request should receive as a response a MessageBird message object like the one below:
```
[{"Id":"5e404e4fef8249da98bdd5cfcc0ce2d6","HRef":"https://rest.messagebird.com/messages/5e404e4fef8249da98bdd5cfcc0ce2d6","Direction":"mt","Type":"sms","Originator":"YourName","Body":"Test message","Reference":"","Validity":null,"Gateway":10,"TypeDetails":{},"DataCoding":"plain","MClass":1,"ScheduledDatetime":null,"CreatedDatetime":"2017-09-17T14:42:30Z","Recipients":{"TotalCount":1,"TotalSentCount":1,"TotalDeliveredCount":0,"TotalDeliveryFailedCount":0,"Items":[{"Recipient":31612345678,"Status":"sent","StatusDatetime":"2017-09-17T14:42:30Z"}]},"Errors":null}]
```

## 3. Implementation notes

### 3.1. Characters and message length

Considering the concatenated SMS limitations and differences between GSM 03.38 characters and unicode characters, the number of messages sent follows the table described [in this article](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message-).

For more about GSM 03.38 characters and which ones from this set counts as 2 characters in a message, check out [this article](https://support.messagebird.com/hc/en-us/articles/208739765-Which-special-characters-count-as-two-characters-in-a-text-message-).

## 4. Known issues

* **Concatenated unicode messages:** Even though the SMS UDH seems correct and the message parts respect the maximum size, the concatenated messages are not showing as a single message in the phone when the characters are unicode.