package dal

import (
	"testing"

	// external
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSapUpload(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SapUpload Suite")
}
