package handler_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/countingtoten/shorty/handler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var (
		h  *handler.Handler
		rw *httptest.ResponseRecorder
		r  *http.Request
	)

	BeforeEach(func() {
		h = handler.New()
		rw = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		h.ServeHTTP(rw, r)
	})

	Describe("POST /new", func() {

	})

	Describe("GET /shorturl", func() {

	})
})
