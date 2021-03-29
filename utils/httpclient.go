package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/sirupsen/logrus"
)

type OptionFunc func(*Options)

func WithTimeout(v time.Duration) OptionFunc {
	return func(o *Options) {
		o.Timeout = v
	}
}

func WithHeader(v http.Header) OptionFunc {
	return func(o *Options) {
		o.Header = v
	}
}

func WithHeaderSet(k, v string) OptionFunc {
	return func(o *Options) {
		o.Header.Set(k, v)
	}
}

func WithToFile(v string) OptionFunc {
	return func(o *Options) {
		o.ToFile = v
	}
}

func WithOKOnly(v bool) OptionFunc {
	return func(o *Options) {
		o.OKOnly = v
	}
}

// Options ...
type Options struct {
	Timeout time.Duration // http.Client.Timeout, default is 5s
	Header  http.Header   // http.Header
	ToFile  string        // save response body to file
	OKOnly  bool          // return error directly if http.StatusCode != 200, default is true
}

// NewOptions ...
func NewOptions(opts ...OptionFunc) *Options {

	option := &Options{
		Timeout: 5 * time.Second,
		Header:  make(http.Header),
		OKOnly:  true,
	}

	for _, opt := range opts {
		opt(option)
	}

	return option
}

func doRequest(method, URL string, data interface{}, options ...*Options) (bs []byte, code int, err error) {

	defer func() {
		if p := recover(); p != nil {
			logrus.Errorf("doRequest exec panic, message is %s, stack is %s", p, string(debug.Stack()))
			err = fmt.Errorf("doRequest exec panic, message is %s", p)
		}
	}()

	option := NewOptions()
	if len(options) > 0 {
		option = options[0]
	}

	client := http.Client{
		Timeout: option.Timeout,
		Transport: &http.Transport{
			// DisableKeepAlives:  true,
			// DisableCompression: true,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	var body io.Reader
	if data != nil {
		switch v := data.(type) {
		case []byte:
			body = bytes.NewBuffer(v)

		case string:
			body = strings.NewReader(v)

		default:
			bs, err = json.Marshal(data)
			if err != nil {
				err = fmt.Errorf("%s URL(%s) marshal data failed, message is %s, raw data is %#v",
					method, URL, err.Error(), data)
				return
			}

			body = bytes.NewBuffer(bs)
		}
	}

	var request *http.Request
	request, err = http.NewRequest(method, URL, body)
	if err != nil {
		logrus.Errorf("%s URL(%s) new request failed, message is %s", method, URL, err.Error())
		return
	}

	// set request header
	request.Header = option.Header

	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		logrus.Errorf("%s URL(%s) failed, message is %s, raw data is %s", method, URL, err.Error(), PrettyPrint(data))
		return
	}
	defer response.Body.Close()

	code = response.StatusCode
	if response.StatusCode != http.StatusOK && option.OKOnly {

		var message string
		bs, err = ioutil.ReadAll(response.Body)
		if err != nil {
			message = fmt.Sprintf("nil(response body read failed by %s)", err.Error())
		} else {
			message = string(bs)
		}

		// 截取字符串，避免错误信息过长
		if len(message) > 1024 {
			message = message[:1024]
		}

		err = fmt.Errorf("%s URL(%s) return code is %d, raw message is %s", method, URL, response.StatusCode, message)
		return
	}

	if len(option.ToFile) > 0 {

		// os.Create auto truncated if the file already exists
		//
		// if bdfile.FileExists(option.ToFile) {
		// 	if err = os.Remove(option.ToFile); err != nil {
		// 		err = fmt.Errorf("%s URL(%s) remove local file failed, message is %s", method, URL, err.Error())
		// 		return
		// 	}
		// }

		var f *os.File
		f, err = os.Create(option.ToFile)
		if err != nil {
			err = fmt.Errorf("%s URL(%s) create local file(%s) failed, message is %s", method, URL, option.ToFile, err.Error())
			return
		}
		defer f.Close()

		_, err = io.Copy(f, response.Body)
		if err != nil {
			err = fmt.Errorf("%s URL(%s) copy file(%s) failed, message is %s", method, URL, option.ToFile, err.Error())
		}

		return
	}

	bs, err = ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("%s URL(%s) response body read failed, message is %s", method, URL, err.Error())
	}

	return
}

// GetURLWithJSONResult ...
func GetURLWithJSONResult(URL string, result interface{}, options ...*Options) (err error) {

	var body []byte
	body, _, err = doRequest("GET", URL, nil, options...)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		err = fmt.Errorf("GET URL(%s) unmarshal result failed, message is %s, raw data is %s",
			URL, err.Error(), string(body))
	}

	return
}

// POSTJSONData ...
func POSTJSONData(URL string, data interface{}, options ...*Options) (err error) {

	option := NewOptions()
	if len(options) > 0 {
		option = options[0]
	}
	option.Header.Set("Content-Type", binding.MIMEJSON)

	_, _, err = doRequest("POST", URL, data, option)
	return
}

// POSTJSONWithJSONResult ...
func POSTJSONWithJSONResult(URL string, data, result interface{}, options ...*Options) (err error) {

	option := NewOptions()
	if len(options) > 0 {
		option = options[0]
	}
	option.Header.Set("Content-Type", binding.MIMEJSON)

	var body []byte
	body, _, err = doRequest("POST", URL, data, option)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		err = fmt.Errorf("POST URL(%s) unmarshal result failed, message is %s, raw data is %s",
			URL, err.Error(), string(body))
	}

	return
}

// DownloadFile ...
func DownloadFile(URL, location string, options ...*Options) (err error) {

	option := NewOptions()
	if len(options) > 0 {
		option = options[0]
	}
	option.ToFile = location

	_, _, err = doRequest("GET", URL, nil, option)
	return
}
