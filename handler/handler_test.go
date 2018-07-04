package handler_test

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/countingtoten/shorty"
	"github.com/countingtoten/shorty/handler"
	"github.com/countingtoten/shorty/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

func newURLPayload() *bytes.Buffer {
	payload := `
		{
		  "user_id": 1,
		  "url": "https://example.com"
		}
	`
	return bytes.NewBuffer([]byte(payload))
}

var emptyPayload = bytes.NewBuffer(nil)

func shortCodeFromURL(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(u.Path, "/")
}

var _ = Describe("Handler", func() {
	var (
		h         *handler.Handler
		datastore *mocks.Datastore
		rw        *httptest.ResponseRecorder
		r         *http.Request
		logs      *gbytes.Buffer
	)

	BeforeEach(func() {
		datastore = &mocks.Datastore{}

		logs = gbytes.NewBuffer()
		logger := zerolog.New(logs)

		h = handler.New(datastore, logger)

		rw = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		h.ServeHTTP(rw, r)
	})

	Describe("/new", func() {
		Context("on success", func() {
			var (
				userID   shorty.UserID
				shortURL shorty.ShortURL = "https://localhost/abcdef"
				longURL  shorty.LongURL
			)

			BeforeEach(func() {
				datastore.On("CreateShortURL", mock.Anything, mock.Anything).
					Return(shortURL, nil).
					Run(func(args mock.Arguments) {
						userID = args.Get(0).(shorty.UserID)
						longURL = args.Get(1).(shorty.LongURL)
					})
				r, _ = http.NewRequest(http.MethodPost, "/new", newURLPayload())
			})

			It("returns a 201 created", func() {
				Expect(rw.Code).To(Equal(http.StatusCreated))
			})

			It("returns the short url in the json payload", func() {
				format := `{ "short_url": "%s" }`
				json := fmt.Sprintf(format, shortURL)

				Expect(rw.Body).To(MatchJSON(json))
			})

			It("calls CreateShortURL with the user id", func() {
				Expect(userID).To(Equal(shorty.UserID(1)))
			})

			It("calls CreateShortURL with the url", func() {
				Expect(longURL).To(Equal("https://example.com"))
			})
		})

		Context("with an incorrect http method", func() {
			BeforeEach(func() {
				r, _ = http.NewRequest(http.MethodGet, "/new", nil)
			})

			It("returns a 404 not found", func() {
				Expect(rw.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with an invalid request length", func() {
			BeforeEach(func() {
				r, _ = http.NewRequest(http.MethodPost, "/new", newURLPayload())
				r.ContentLength = 1000
			})

			It("returns a 400 bad request", func() {
				Expect(rw.Code).To(Equal(http.StatusBadRequest))
			})

			It("logs an error message", func() {
				Eventually(logs).Should(gbytes.Say(`Handler.CreateShortURL unable to read request body`))
			})
		})

		Context("with an invalid request length", func() {
			BeforeEach(func() {
				r, _ = http.NewRequest(http.MethodPost, "/new", emptyPayload)
			})

			It("returns a 400 bad request", func() {
				Expect(rw.Code).To(Equal(http.StatusBadRequest))
			})

			It("logs an error message", func() {
				Eventually(logs).Should(gbytes.Say(`Handler.CreateShortURL unable to parse request`))
			})
		})

		Context("when CreateShortURL returns an error", func() {
			var (
				err error
			)

			BeforeEach(func() {
				err = errors.New("CreateShortURL failed")
				datastore.On("CreateShortURL", mock.Anything, mock.Anything).Return("", err)
				r, _ = http.NewRequest(http.MethodPost, "/new", newURLPayload())
			})

			It("returns a 500 server error", func() {
				Expect(rw.Code).To(Equal(http.StatusInternalServerError))
			})

			It("logs an error message", func() {
				Eventually(logs).Should(gbytes.Say(`Handler.CreateShortURL unable to create short url`))
			})

			It("includes the error message in the logs", func() {
				Eventually(logs).Should(gbytes.Say(err.Error()))
			})
		})
	})

	Describe("/shorturl", func() {
		var (
			shortURL shorty.ShortURL = "https://localhost/abcdef"
			longURL  shorty.LongURL  = "https://longurl.com/some/path"
			arg      string
		)

		BeforeEach(func() {
			r, _ = http.NewRequest(http.MethodGet, shortURL, nil)
		})

		Context("on success", func() {
			BeforeEach(func() {
				datastore.On("GetLongURL", mock.Anything).
					Return(longURL, nil).
					Run(func(args mock.Arguments) {
						arg = args.Get(0).(string)
					})
			})

			It("redirects to the long url", func() {
				Expect(rw.Header().Get("Location")).To(Equal(longURL))
			})

			It("sets the status code to 301 permanently moved", func() {
				Expect(rw.Code).To(Equal(http.StatusMovedPermanently))
			})

			It("calls GetLongURL with the request url as an argument", func() {
				shortCode := shortCodeFromURL(shortURL)
				Expect(arg).To(Equal(shortCode))
			})
		})

		Context("when GetLongURL returns a blank LongURL", func() {
			BeforeEach(func() {
				datastore.On("GetLongURL", mock.Anything).Return("", nil)
			})

			It("sets the status code to 404 not found", func() {
				Expect(rw.Code).To(Equal(http.StatusNotFound))
			})

			It("logs a warning message about the missing path", func() {
				Eventually(logs).Should(gbytes.Say(`Handler.GetURL 404`))
			})
		})

		Context("when GetLongURL returns an error", func() {
			var (
				err error
			)

			BeforeEach(func() {
				err = errors.New("GetLongURL failed")
				datastore.On("GetLongURL", mock.Anything).Return("", err)
			})

			It("sets the status code to 500 internal server error", func() {
				Expect(rw.Code).To(Equal(http.StatusInternalServerError))
			})

			It("logs an error message", func() {
				Eventually(logs).Should(gbytes.Say(`Handler.GetURL unable to get the long url`))
			})

			It("includes the error message in the logs", func() {
				Eventually(logs).Should(gbytes.Say(err.Error()))
			})
		})
	})
})
