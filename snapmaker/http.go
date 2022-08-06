package snapmaker

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpRequestMatcher struct {
	expectedRequest *http.Request
}

func NewHttpRequestMatcher(req *http.Request) *HttpRequestMatcher {
	return &HttpRequestMatcher{expectedRequest: req}
}

func (m *HttpRequestMatcher) Matches(x interface{}) bool {
	actualRequest, isHttpRequest := x.(*http.Request)
	switch {
	case !isHttpRequest:
		return false
	case !reflect.DeepEqual(actualRequest.URL, m.expectedRequest.URL):
		return false
	case !reflect.DeepEqual(actualRequest.Header, m.expectedRequest.Header):
		return false
	}
	return true
}

func (m *HttpRequestMatcher) String() string {
	return fmt.Sprintf("%v", m.expectedRequest)
}

func isEqual(a, b *url.URL) bool {

	switch {
	case a.RawPath != b.RawPath:
		return false
	case a.Host != b.Host:
		return false
	case a.Scheme != b.Scheme:
		return false
	case !reflect.DeepEqual(a.Query().Get("token"), b.Query().Get("token")):
		// only compare token URL parameter, the timestamp parameter is too volatile
		return false
	}

	return true
}
