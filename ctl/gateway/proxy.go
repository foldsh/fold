package gateway

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Serve() {
	r := gin.Default()
	r.Any("/:service/*path", func(c *gin.Context) {
		proxy(c)
	})
	r.Run(":6123")
}

func proxy(c *gin.Context) error {
	service := c.Param("service")
	urlStr := fmt.Sprintf("http://%s:6123", service)
	log.Println(urlStr)
	remote, err := url.Parse(urlStr)
	if err != nil {
		return err
	}
	// Yes, this sets up the proxy all over again per request... but it works fine
	// for this local development use case.
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Method = c.Request.Method
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = c.Param("path")
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		if strings.Contains(err.Error(), "no such host") {
			// This means that the the host wasn't found and therefore that the service name
			// given in the url is incorrect.
			// It's rubbish we have to match this error using a string search but the 'no such host'
			// error is not exported:
			// https://golang.org/src/net/net.go?h=no+such+host
			// Perhaps because of this, I can't seem to get it work with errors.Is
			w.WriteHeader(404)
			w.Write(
				[]byte(
					fmt.Sprintf(
						`{"title":"%s","service-name":"%s","detail":"%s"}`,
						serviceNotFound,
						service,
						snfDetail,
					),
				),
			)
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
	return nil
}

var (
	serviceNotFound = "Service not found"
	snfDetail       = "The service you tried to send a request to does not exist. Either you have spelled it incorrectly or it is not currently running."
)
