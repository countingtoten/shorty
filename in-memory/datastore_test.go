package memory_test

import (
	"math/rand"
	"net/url"
	"strings"

	"github.com/caarlos0/env"
	"github.com/countingtoten/shorty"
	memory "github.com/countingtoten/shorty/in-memory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func shortCodeFromURL(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}

	return strings.TrimPrefix(u.Path, "/")
}

var _ = Describe("Datastore", func() {
	var (
		ds  *memory.Datastore
		cfg shorty.Config
	)

	BeforeEach(func() {
		cfg = shorty.Config{}
		env.Parse(&cfg)

		ds = memory.NewDatastore(cfg)
	})

	Describe("CreateShortURL", func() {
		var (
			id       = rand.Int63()
			url      = "https://example.com"
			shortURL shorty.ShortURL
			err      error
		)

		BeforeEach(func() {
			shortURL, err = ds.CreateShortURL(id, url)
		})

		It("creates a user if id not found", func() {
			Expect(ds.UserData).To(HaveKey(id))
		})

		It("creates a value on ShortURLs", func() {
			shortCode := shortCodeFromURL(shortURL)
			Expect(ds.ShortCodes).To(HaveKey(shortCode))
		})

		It("includes the base url", func() {
			Expect(shortURL).To(ContainSubstring(cfg.BaseURL))
		})

		It("only contains alphanumeric characters in the path", func() {
			shortCode := shortCodeFromURL(shortURL)
			Expect(shortCode).To(MatchRegexp(`^[A-Za-z0-9]+$`))
		})

		It("does not return an error", func() {
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetLongURL", func() {
		var (
			id       = rand.Int63()
			url      = "https://example.com"
			shortURL shorty.ShortURL
			longURL  shorty.LongURL
			err      error
		)

		Context("on success", func() {
			BeforeEach(func() {
				shortURL, err = ds.CreateShortURL(id, url)

				shortCode := shortCodeFromURL(shortURL)

				longURL, err = ds.GetLongURL(shortCode)
			})

			It("returns the LongURL associated with the ShortURL", func() {
				Expect(longURL).To(Equal(url))
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with a ShortURL parameter that does not exist", func() {
			BeforeEach(func() {
				longURL, err = ds.GetLongURL("")
			})

			It("returns an empty string", func() {
				Expect(longURL).To(BeEmpty())
			})

			It("does not return an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("GetUser", func() {
		var (
			id = rand.Int63()
		)

		BeforeEach(func() {
			ds.GetUser(id)
		})

		It("creates a user if id not found", func() {
			Expect(ds.UserData).To(HaveKey(id))
		})
	})

	Describe("NewShortURL", func() {
		var (
			shortCode shorty.ShortCode
		)

		BeforeEach(func() {
			shortCode = ds.NewShortCode()
		})

		It("only contains alphanumeric characters in the path", func() {
			Expect(shortCode).To(MatchRegexp(`^[A-Za-z0-9]+$`))
		})
	})
})
