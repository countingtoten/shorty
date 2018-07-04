package memory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "InMemory Suite")
}
