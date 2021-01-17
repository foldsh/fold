// Package handler defines a variety of handlers which are used to implement
// runtime interfaces for different environments. A handlers job is basically
// to take input from the environment, and translate it into the appropriate
// action for a service.
//
// For example, an HTTP handler will receive an incoming HTTP request and use
// it to form a `Request` and then call `DoRequest` on the `Service`. It will
// also take the `Response` and form an appropriate HTTP response. The lambda
// handler does the same thing, but translates to/from events passed to the
// fold runtime by AWS Lambda.
package handler

type Handler interface {
}
