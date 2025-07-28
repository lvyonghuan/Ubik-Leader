package call

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/Ubik-Util/ujson"
	"github.com/lvyonghuan/Ubik-Util/ulog"
)

type Caller struct {
	log ulog.Log

	Followers map[string]Follower // Map of follower UUIDs to their addresses
}

type response struct {
	Info   []byte `json:"info"`
	Status int    `json:"status"`
}

func InitCaller(ulog ulog.Log) *Caller {
	caller := Caller{
		log:       ulog,
		Followers: make(map[string]Follower),
	}

	caller.log.Debug("Caller initialized")

	return &caller
}

// CallAndUnmarshal calls the specified URL with the given method and body,
// and unmarshal the response into the provided variable.
func (c *Caller) callAndUnmarshal(method, url string, body, v any) (int, error) {
	c.log.Debug("Calling URL: " + url)

	// Create a new request
	req, err := c.newRequest(method, url, body)
	if err != nil {
		return 0, uerr.NewError(err)
	}

	// Call the request
	resp, err := c.call(req)
	if err != nil {
		return 0, uerr.NewError(err)
	}

	// Preliminary unmarshal of the response
	res, err := c.preliminaryUnmarshalResponse(resp)
	if err != nil {
		return 0, uerr.NewError(err)
	}

	// Unmarshal Info field
	if v != nil {
		err = ujson.Unmarshal(res.Info, v)
	}

	return res.Status, nil
}

func (c *Caller) newRequest(method, url string, body any) (*http.Request, error) {
	var bodyReader = new(bytes.Buffer)
	if body != nil {
		bodyReader = bytes.NewBuffer(body.([]byte))
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, uerr.NewError(err)
	}

	return req, nil
}

func (c *Caller) call(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, uerr.NewError(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, uerr.NewError(errors.New("request failed with status code: " + strconv.Itoa(resp.StatusCode)))
	}

	return resp, nil
}

func (c *Caller) preliminaryUnmarshalResponse(resp *http.Response) (response, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response{}, uerr.NewError(errors.New("response status not OK: " + strconv.Itoa(resp.StatusCode)))
	}

	var res response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response{}, uerr.NewError(err)
	}
	if err := ujson.Unmarshal(body, &res); err != nil {
		return response{}, err
	}

	return res, nil
}
