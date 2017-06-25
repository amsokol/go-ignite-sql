package http

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Client contains Ignite cluster connection information
type Client struct {
	ConnectionInfo *ConnectionInfo
}

func (c *Client) execute(v *url.Values) ([]byte, error) {
	// TODO: add round-robin to select node
	req, err := http.NewRequest("POST", c.ConnectionInfo.Servers[0], strings.NewReader(v.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "Can't create new POST http.Request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Can't Do HTTP request by DefaultClient")
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Can't read bytes from HTTP response body")
	}

	return b, err
}

func (c *Client) getError(code int, str string) string {
	if code < successStatusSuccess || successStatusSecurityCheckFailed < code {
		code = successStatusUnknown
	}
	m := strings.Join([]string{"Ignite returns: ", successStatusMsg[code]}, "")
	if len(str) > 0 {
		m = strings.Join([]string{m, str}, ": ")
	}
	return m
}