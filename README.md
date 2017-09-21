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

Go to the project root folder and fill the configurations in *config.json* file. This file should looks like the example below:
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

* **MbApiKey:** Your MessageBird REST API key. The *development* entry corresponds to the *Mode test* key and the *production key* to the *Mode live* key. Notice that it's not required to setup both of them, but the application will not work properly if it's running in a environment with the correspondent key not set.
* **MaxCsmsParts:** Maximum number of concatenated SMS parts. Following MessageBird Dashboard Quick Send interface, it was initially set to 9 but it is not a mandatory parameter and, if removed from the configuration, it **will be set to 255** automatically due to the limitation of number of parts that can be encoded in a CSMS UDH. Must be **between 0 and 255**.

At last, you can check if your configuration file is fine by running the command `go test -run TestLoadConfig`.

### 1.3. Build and Run

Build the project by typing the command `go build`, which will generate the binary file `mbassignment`.

To run the API, simply type `./mbassignment` in your terminal. If the project is running properly, the API will be up and running at `localhost:8080` and you will see a message like this one:
```
2017/09/17 10:56:01 ENV: development
2017/09/17 10:56:01 Server started!
```

The server will run by default in the development environment and thus use MessageBird API test key. To use the live key, please run the API with the *production environment variable*:
```
ENV=production ./mbassignment
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
A successful request should receive as a response a [MessageBird message object](https://developers.messagebird.com/docs/messaging#messaging-object) like the one below:
```
[{"Id":"5e404e4fef8249da98bdd5cfcc0ce2d6","HRef":"https://rest.messagebird.com/messages/5e404e4fef8249da98bdd5cfcc0ce2d6","Direction":"mt","Type":"sms","Originator":"YourName","Body":"Test message","Reference":"","Validity":null,"Gateway":10,"TypeDetails":{},"DataCoding":"plain","MClass":1,"ScheduledDatetime":null,"CreatedDatetime":"2017-09-17T14:42:30Z","Recipients":{"TotalCount":1,"TotalSentCount":1,"TotalDeliveredCount":0,"TotalDeliveryFailedCount":0,"Items":[{"Recipient":31612345678,"Status":"sent","StatusDatetime":"2017-09-17T14:42:30Z"}]},"Errors":null}]
```

## 3. Implementation notes

### 3.1. Files and modules

Here's a quick description of the project files/modules:

* **main.go:** Instantiates and runs the HTTP server.
* **config.go:** Loads the configurations in *config.json*.
* **router.go:** Defines the application routes and their correspondent handlers.
* **handlers.go:** Defines the methods which will handle the server requests and their responses.
* **validator.go:** Validates the request parameters.
* **mbapi.go:** A wrapper for MessageBird API client which deals with the message sending and the API throughput.
* **text_helper:** A helper to find out the message *datacoding* and to split the message body when necessary with its proper length.
* **\*_test.go:** Unit tests for its corresponding module.

### 3.2. Datacoding and message length

Considering MessageBird SMS datacoding types - plain and unicode, the SMS body was classified as plain if all its characters were contained in [this GSM 03.38 table](https://en.wikipedia.org/wiki/GSM_03.38#GSM_7-bit_default_alphabet_and_extension_table_of_3GPP_TS_23.038_.2F_GSM_03.38) and unicode if at least one character was not contained. Moreover, the number of messages sent according to the total number of characters and datacoding follows the table described [in this article](https://support.messagebird.com/hc/en-us/articles/208739745-How-long-is-1-SMS-Message-). 

At last, the following characters were considered as 2 characters in a plain message: `\n`, `\`, `^`, `~`, `[`, `]`, `{`, `}`, `|`, `~`, `€`.

### 3.3. UDH Reference Number

As described in [this article in Wikipedia](https://en.wikipedia.org/wiki/Concatenated_SMS), the UDH reference number must be the same for all SMS parts in the CSMS and can be composed by either 1 or 2 octets. For the sake of simplicity, this assignment uses a 8-bit reference number.

The reference number generation was implemented as a sequential number starting as a random number in the interval [0, 256) in order to avoid repeated numbers when restarting the application. This number is incremented for every new group of CSMSes sent and is reset to zero when its value is greater than its maximum value (255).

**Note:** I have considered that a more sophisticated solution would make the assignment more complex while not avoiding the possibility of reference number collision since: **(1)** there is no way to guarantee a SMS delivery time (or if it will be in fact delivered) and delivery order; and **(2)** in practice the UDH is not the only identifying factor in sending an SMS.

## 4. Known issues

* **Unicode characters:** Sometimes unicode characters are not properly rendered, which seems to be related to the mobile network carrier capacity to transcre it, as reported in other platforms such as [Elvanto](https://help.elvanto.com/using-elvanto/emails-letters-sms/sms/information-about-sms-character-encoding-language-support-and-message-lengths/) and [CRM Text](http://crmtext.com/api/sms-api-understandingcarriers).
* **Binary messages:** I've experienced the issue of binary messages being displayed with question marks in it's the beginning and concatenated SMSes displays as separated messages. This can be related to the previous issue or may be related to [the issue described here](https://issuetracker.google.com/issues/36944392#c5). If you run into this problem, please check if the same issue can be replicated from MessageBird Dashboard Quick Send and also try sending the SMS to another phone.
* **API throughput in a Virtual Machine (VM):** The API throughput was implemented using Golang time.Ticker and tested in a Lubuntu VM, where I have noticed that the timer was ticking a few milliseconds short to its due time. As discussed in [this thread](https://github.com/golang/go/issues/19810), this seems to not be related to Golang, but to the VM monotonic clock not being synchronised with the real time clock.