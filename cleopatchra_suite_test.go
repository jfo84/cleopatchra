package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCleopatchra(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cleopatchra Suite")
}
