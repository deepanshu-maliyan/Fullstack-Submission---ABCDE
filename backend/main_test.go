package main_test

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEcommerce(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ecommerce Suite")
}
