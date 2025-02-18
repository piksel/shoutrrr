package telegram_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/containrrr/shoutrrr/pkg/plugins/telegram"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTelegram(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shoutrrr Telegram Suite")
}

var _ = Describe("the telegram plugin", func() {
	var telegram *Plugin
	var envTelegramUrl string

	BeforeSuite(func() {
		telegram = &Plugin{}
		envTelegramUrl = os.Getenv("SHOUTRRR_TELEGRAM_URL")

	})


	When("running integration tests", func() {
		It("should not error out", func() {
			if envTelegramUrl == "" {
				return
			}
			err := telegram.Send(envTelegramUrl, "This is an integration test message")
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("sending a message", func() {
		When("given a valid request with a faked token", func() {
			It("should generate a 401", func() {
				url := "telegram://000000000:AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA/channel-id"
				message := "this is a perfectly valid message"
				err := telegram.Send(url, message)
				Expect(err).To(HaveOccurred())
				fmt.Println(err.Error())
				Expect(strings.Contains(err.Error(), "401 Unauthorized")).To(BeTrue())
			})
		})
		When("given a message that exceeds the max length", func() {
			It("should generate an error", func() {
				hundredChars := "this string is exactly (to the letter) a hundred characters long which will make the send func error"
				url := "telegram://12345:mock-token/channel-1"
				builder := strings.Builder{}
				for i := 0; i < 42; i++ {
					builder.WriteString(hundredChars)
				}

				err := telegram.Send(url, builder.String())
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("creating configurations", func() {
		When("given an url", func() {
			It("should return an error if no arguments where supplied", func() {
				url := "telegram://"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should return an error if the token has an invalid format", func() {
				url := "telegram://invalid-token"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should return an error if only the api token where supplied", func() {
				url := "telegram://12345:mock-token"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).To(HaveOccurred())
				Expect(config == nil).To(BeTrue())
			})
			It("should create a config object", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config != nil).To(BeTrue())
			})
			It("should create a config object containing the API Token", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.Token).To(Equal("12345:mock-token"))
			})
			It("should add every subsequent argument as a channel id", func() {
				url := "telegram://12345:mock-token/channel-1/channel-2/channel-3"
				config, err := telegram.CreateConfigFromURL(url)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.Channels).To(Equal([]string {
					"channel-1",
					"channel-2",
					"channel-3",
				}))
			})
		})
	})
})

