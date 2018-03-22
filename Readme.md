# Health Checker

It is a simple tool to monitor your web service and notify you once something go wrong.

## Usage

```bash
Usage of hc:
  -config string
        path of config.json
```

You can start hc like this:
```bash
hc --config /tmp/config.json
```
The default dir for config.json is ```/home/xxx/.hc/config.json```

## Config

> Notice: default checker will check the status of response, make sure your service response 200.

Here is a simple example:
```json
{
  "name": "HealthChecker",
  "url": "https://github.com/sleagon",
  "mail": {
    "username": "shanyy163@163.com",
    "password": "*******",
    "host": "smtp.163.com",
    "port": 25
  },
  "plans": {
    "bad": {
      "url": "https://www.bgaidu.com",
      "method": "HEAD",
      "body": "{}",
      "interval": 5,
      "callback": "https://www.smartstudy.com",
      "mail": "shanyuanyuan@innobuddy.com"
    },
    "good": {
      "url": "https://www.baidu.com",
      "method": "HEAD",
      "body": "{}",
      "interval": 5,
      "callback": "https://www.smartstudy.com",
      "mail": "shanyuanyuan@innobuddy.com"
    },
    "best": {
      "url": "https://www.google.com",
      "method": "HEAD",
      "body": "{}",
      "interval": 5,
      "callback": "https://www.smartstudy.com",
      "mail": "shanyuanyuan@innobuddy.com"
    }
  }
}
```

The ```config.mail``` is used for sending notice. You may need set it according to your email address.

For gmail: https://support.google.com/a/answer/176600?hl=en

For 163: http://help.163.com/09/1223/14/5R7P3QI100753VB8.html

> Gmail's port 465 is for connecting via TLS, but SendMail expects plain old TCP. Try connecting to port 587 instead. SendMail will upgrade to TLS automatically when it's available (which it is in this case). You should used the port, like 25, which do not require ssl.

Just add the service you want observe in ```config.plans```.

|key|usage|
|-|-|
|url|Your service url|
|method|Method used for request, HEAD is recommended. You can GET/POST too.|
|body|Body you want to include in the request|
|interval|In second|
|callback|The service you want to call when service is offline, used for reload your service. Just like the hook in github|
|mail| Email address used for receiving notice when your sevice is down.|



