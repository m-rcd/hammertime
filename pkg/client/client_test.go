package client_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/flintlock/api/services/microvm/v1alpha1"
	"github.com/weaveworks/flintlock/api/types"
	"github.com/weaveworks/flintlock/client/cloudinit"
	"gopkg.in/yaml.v2"

	"github.com/Callisto13/hammertime/pkg/client"
	"github.com/Callisto13/hammertime/pkg/client/fakeclient"
)

var _ = Describe("Client", func() {
	var (
		name      = "Pantalaimon"
		namespace = "Casper"
	)

	It("creates a MicroVm", func() {
		mockClient := new(fakeclient.FakeMicroVMClient)
		c := client.New(mockClient)

		metaData, err := yaml.Marshal(cloudinit.Metadata{
			InstanceID:    fmt.Sprintf("%s/%s", namespace, name),
			LocalHostname: name,
			Platform:      "liquid_metal",
		})
		Expect(err).ToNot(HaveOccurred())
		metadata := base64.StdEncoding.EncodeToString(metaData)

		userData := &cloudinit.UserData{
			HostName: name,
			Users: []cloudinit.User{
				{
					Name: "root",
				},
			},
			FinalMessage: "The Liquid Metal booted system is good to go after $UPTIME seconds",
			BootCommands: []string{
				"ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf",
			},
		}
		data, err := yaml.Marshal(userData)
		Expect(err).ToNot(HaveOccurred())

		dataWithHeader := append([]byte("#cloud-config\n"), data...)
		userdata := base64.StdEncoding.EncodeToString(dataWithHeader)
		microVm := &v1alpha1.CreateMicroVMResponse{
			Microvm: &types.MicroVM{
				Spec: &types.MicroVMSpec{
					Id:        name,
					Namespace: namespace,
				},
			},
		}

		mockClient.CreateMicroVMReturns(microVm, nil)
		result, err := c.Create(name, namespace, "", "")
		Expect(err).ToNot(HaveOccurred())

		_, inputReq, _ := mockClient.CreateMicroVMArgsForCall(0)
		Expect(inputReq.Microvm.Id).To(Equal(name))
		Expect(inputReq.Microvm.Namespace).To(Equal(namespace))
		Expect(inputReq.Microvm.Metadata["meta-data"]).To(Equal(metadata))
		Expect(inputReq.Microvm.Metadata["user-data"]).To(Equal(userdata))

		Expect(mockClient.CreateMicroVMCallCount()).To(Equal(1))
		Expect(result.Microvm.Spec.Id).To(Equal(name))
		Expect(result.Microvm.Spec.Namespace).To(Equal(namespace))
	})

	Context("when using sshkey", func() {
		Context("when file exists", func() {
			var (
				keyFile *os.File
				key     = "this is a test key woohoo"
			)

			BeforeEach(func() {
				var err error
				keyFile, err = ioutil.TempFile("", "ssh_key")
				Expect(err).NotTo(HaveOccurred())

				err = ioutil.WriteFile(keyFile.Name(), []byte(key), 0)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				Expect(os.Remove(keyFile.Name())).To(Succeed())
			})

			It("creates a MicroVm", func() {
				mockClient := new(fakeclient.FakeMicroVMClient)
				c := client.New(mockClient)

				metaData, err := yaml.Marshal(cloudinit.Metadata{
					InstanceID:    fmt.Sprintf("%s/%s", namespace, name),
					LocalHostname: name,
					Platform:      "liquid_metal",
				})
				Expect(err).ToNot(HaveOccurred())
				metadata := base64.StdEncoding.EncodeToString(metaData)

				userData := &cloudinit.UserData{
					HostName: name,
					Users: []cloudinit.User{
						{
							Name:              "root",
							SSHAuthorizedKeys: []string{key},
						},
					},
					FinalMessage: "The Liquid Metal booted system is good to go after $UPTIME seconds",
					BootCommands: []string{
						"ln -sf /run/systemd/resolve/stub-resolv.conf /etc/resolv.conf",
					},
				}
				data, err := yaml.Marshal(userData)
				Expect(err).ToNot(HaveOccurred())

				dataWithHeader := append([]byte("#cloud-config\n"), data...)
				userdata := base64.StdEncoding.EncodeToString(dataWithHeader)
				microVm := &v1alpha1.CreateMicroVMResponse{
					Microvm: &types.MicroVM{
						Spec: &types.MicroVMSpec{
							Id:        name,
							Namespace: namespace,
						},
					},
				}

				mockClient.CreateMicroVMReturns(microVm, nil)
				result, err := c.Create(name, namespace, "", keyFile.Name())
				Expect(err).ToNot(HaveOccurred())

				_, inputReq, _ := mockClient.CreateMicroVMArgsForCall(0)
				Expect(inputReq.Microvm.Id).To(Equal(name))
				Expect(inputReq.Microvm.Namespace).To(Equal(namespace))
				Expect(inputReq.Microvm.Metadata["meta-data"]).To(Equal(metadata))
				Expect(inputReq.Microvm.Metadata["user-data"]).To(Equal(userdata))

				Expect(mockClient.CreateMicroVMCallCount()).To(Equal(1))
				Expect(result.Microvm.Spec.Id).To(Equal(name))
				Expect(result.Microvm.Spec.Namespace).To(Equal(namespace))
			})
		})

		Context("when file does not exist", func() {
			It("returns an error", func() {
				mockClient := new(fakeclient.FakeMicroVMClient)
				c := client.New(mockClient)
				_, err := c.Create(name, namespace, "", "key.txt")
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})

	Context("jsonSpec is set", func() {
		var (
			jsonSpec  = "./../../example.json"
			name      = "mvm1"
			namespace = "ns1"
		)

		It("creates a MicroVm", func() {
			mockClient := new(fakeclient.FakeMicroVMClient)
			c := client.New(mockClient)

			microVm := &v1alpha1.CreateMicroVMResponse{
				Microvm: &types.MicroVM{
					Spec: &types.MicroVMSpec{
						Id:        name,
						Namespace: namespace,
					},
				},
			}

			mockClient.CreateMicroVMReturns(microVm, nil)
			result, err := c.Create("", "", jsonSpec, "")
			Expect(err).ToNot(HaveOccurred())

			_, inputReq, _ := mockClient.CreateMicroVMArgsForCall(0)
			Expect(inputReq.Microvm.Id).To(Equal(name))
			Expect(inputReq.Microvm.Namespace).To(Equal(namespace))

			Expect(mockClient.CreateMicroVMCallCount()).To(Equal(1))
			Expect(result.Microvm.Spec.Id).To(Equal(name))
			Expect(result.Microvm.Spec.Namespace).To(Equal(namespace))
		})

		Context("when file does not exist", func() {
			It("returns an error", func() {
				mockClient := new(fakeclient.FakeMicroVMClient)
				c := client.New(mockClient)
				_, err := c.Create("", "", "./../../example1.json", "")
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})
	})
})
