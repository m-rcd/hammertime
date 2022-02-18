package client_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/Callisto13/hammertime/pkg/client"
	"github.com/Callisto13/hammertime/pkg/client/fakeclient"
)

var _ = Describe("Client", func() {
	It("creates a MicroVm", func() {
		mockClient := new(fakeclient.FakeMicroVMClient)
		c := client.New(mockClient)
		_, err := c.Create()
		Expect(err).ToNot(HaveOccurred())
	})
})
