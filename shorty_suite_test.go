package shorty_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestShorty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shorty Suite")
}
