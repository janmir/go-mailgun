package mailgun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

/*
Request
> POST /v3/DOMAIN/messages HTTP/1.1
> Host: api.mailgun.net
> Authorization: Basic api:***REDUCTED***
> User-Agent: insomnia/5.16.6
> Content-Type: application/x-www-form-urlencoded
> Accept:
> Content-Length: 76
> from=jp&to=email0gmail.com&subject=Hello&text=World!

Response
{
	"id": "<20180722005422.1.AFD44FBD540C7AFB@Domain>",
	"message": "Queued. Thank you."
}
*/

const (
	_maingunDomain = "https://api.mailgun.net/v3/%s/messages"
	_timeout       = 5 //Seconds
	_defaultUser   = "api"
)

var (
	httpTimeout time.Duration
	debug       = false
)

//Mail this defines the email
type Mail struct {
	api                     string
	domain                  string
	client                  http.Client
	To, From, Subject, Text string
}

//Response received response
type Response struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func init() {
	httpTimeout = time.Duration(_timeout * time.Second)
}

//DebugMailClient creates the default email client
func DebugMailClient(domain, api string) Mail {
	debug = true
	return Mail{
		api:    api,
		domain: domain,
		client: http.Client{
			Timeout:       httpTimeout,
			CheckRedirect: redirectPolicyFunc,
			Transport:     NewLoggedTransport(http.DefaultTransport, newLogger()),
		},
	}
}

//DefaultMailClient creates the default email client
func DefaultMailClient(domain, api string) Mail {
	return Mail{
		api:    api,
		domain: domain,
		client: http.Client{
			Timeout:       httpTimeout,
			CheckRedirect: redirectPolicyFunc,
		},
	}
}

//Create creates email
func (mail *Mail) Create(to, from, subject, text string) {
	mail.To = to
	mail.From = from
	mail.Subject = subject
	mail.Text = text
}

//Send sends the mail
func (mail *Mail) Send() (Response, error) {
	out := Response{}

	data := url.Values{}
	data.Set("to", mail.To)
	data.Set("from", mail.From)
	data.Set("subject", mail.Subject)
	data.Set("text", mail.Text)

	url := fmt.Sprintf(_maingunDomain, mail.domain)
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return out, fmt.Errorf("Error occured while creating the mail request, %#v", err)
	}
	req.SetBasicAuth(_defaultUser, mail.api)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	catch(httputil.DumpRequestOut(req, true))

	res, err := mail.client.Do(req)
	if err != nil {
		return out, fmt.Errorf("Error occured while performing the mail send, %#v", err)
	}
	defer res.Body.Close()
	catch(httputil.DumpResponse(res, true))

	js, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return out, fmt.Errorf("Error occured while reading mail body, %#v", err)
	}

	err = json.Unmarshal(js, &out)
	if err != nil {
		return out, fmt.Errorf("Error occured while parsinf response, %#v", err)
	}

	return out, nil
}

func catch(data []byte, err error) {
	if !debug {
		return
	}
	if err == nil {
		fmt.Printf("\n%s\n", data)
	} else {
		log.Fatalf("\n%s\n", err)
	}
}
