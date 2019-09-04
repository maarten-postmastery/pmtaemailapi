# PowerMTA Email API

Submit emails to PowerMTA with HTTP and JSON instead of SMTP and MIME.

The Email API is used to submit email messages for delivery. You specify address headers, subject, text, and html. The Email API assembles the email, encodes it according to the appropriate standards, and submits it to PowerMTA.

## Getting started

Start on PowerMTA server with:

    pmtaemailapi -listen 127.0.0.1:8000 -pickup /var/pickup

Submit email messages by posting JSON to http://127.0.0.1:8000/messages.

    curl http://127.0.0.1:8000/messages -d '{"from":"postmaster@sender.com","to":"nobody@example.com","subject":"Test","text":"This is a test"}'

## API reference

The JSON message object can contain the following name/value pairs:

|Name   |Type               |Description                    |
|-------|-------------------|-------------------------------|
|sender |string             |Envelope sender address        |
|from   |object/string      |Address to use in from header  |
|to     |array/object/string|Address(es) to use in to header|
|cc     |array/object/string|Address(es) to use in cc header|
|bcc    |array/object/string|Address(es) to use as bcc      |
|subject|string             |Subject line                   |
|text   |string             |Content for text part          |
|html   |string             |Content for html part          |

An address can be specified as object or as string. An address object contains the following name/value pairs:

|Name   |Type               |Description                    |
|-------|-------------------|-------------------------------|
|name   |string             |Display name                   |
|address|string             |Email address                  |

All strings are expected to be encoded as UTF-8. Display names are MIME word-encoded if needed. Text and HTML bodies are transfer-encoded with quoted-printable.

There is no access control. Make sure that the API is listening on private IPs or use a firewall.

## Example code

### Python

Install requests library with `easy_install requests`.

    import json
    import requests

    url = "http://localhost:8000/messages"
    headers = {'content-type': 'application/json'}
    payload = {
            'from': 'Jim <jim@foo.bar>', 
            'to': 'you@example.com',
            'subject': 'test',
            'text': 'This is a test'}
    r = requests.post(url, data=json.dumps(payload), headers=headers)
    print r.status_code
    print r.text





