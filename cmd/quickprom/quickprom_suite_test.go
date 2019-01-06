package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQuickprom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Quickprom Suite")
}
