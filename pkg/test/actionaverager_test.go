package actionaverager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/action-averager/pkg/actionaverager"
)

const (
	emptyStats = "[]"
)

func addMultipleActions(averager actionaverager.ActionAverager, actions []string) {
	for i := range actions {
		err := averager.AddAction(actions[i])
		Expect(err).NotTo(HaveOccurred())
	}
}

func verifyMultipleDifferentStats(stats string, expSubstrings []string) {
	for i := range expSubstrings {
		//NOTE: use ContainSubstring, since order of actions is go map order so explicit string equals fails intermittently
		Expect(stats).To(ContainSubstring(expSubstrings[i]))
	}
}

var _ = Describe("action-averager tests", func() {
	var averager actionaverager.ActionAverager
	BeforeEach(func() {
		averager = actionaverager.NewActionAverager()
	})

	Context("correct input provided", func() {
		Context("serial function calls", func() {
			It("should not average a single input", func() {
				err := averager.AddAction(`{"action":"run","time":20}`)
				Expect(err).NotTo(HaveOccurred())
				stats := averager.GetStats()
				Expect(stats).To(Equal(`[{"action":"run","avg":20}]`))
			})

			It("should not average multiple inputs to different actions", func() {
				actions := []string{
					`{"action":"run","time":55}`,
					`{"action":"skip","time":145}`,
				}
				addMultipleActions(averager, actions)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"run","avg":55}`,
					`{"action":"skip","avg":145}`,
				}
				verifyMultipleDifferentStats(stats, expStats)
			})

			It("should average multiple inputs of the same action", func() {
				actions := []string{
					`{"action":"hop","time":55.5}`,
					`{"action":"hop","time":145.37}`,
				}
				addMultipleActions(averager, actions)
				stats := averager.GetStats()
				Expect(stats).To(Equal(`[{"action":"hop","avg":100.435}]`))
			})

			It("should average multiple inputs of different actions", func() {
				actions := []string{
					`{"action":"hop","time":55.75}`,
					`{"action":"skip","time":155.123}`,
					`{"action":"jump","time":35.57}`,
					`{"action":"skip","time":155.5}`,
					`{"action":"skip","time":125.545}`,
					`{"action":"jump","time":30}`,
					`{"action":"hop","time":62.3}`,
				}
				addMultipleActions(averager, actions)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"hop","avg":59.025}`,
					`{"action":"skip","avg":145.38933333333333}`,
					`{"action":"jump","avg":32.785}`,
				}
				verifyMultipleDifferentStats(stats, expStats)
			})

			It("should handle a mix of single and multiple inputs to different actions", func() {
				actions := []string{
					`{"action":"walk","time":200}`,
					`{"action":"run","time":100}`,
					`{"action":"crawl","time":300}`,
					`{"action":"run","time":150}`,
					`{"action":"walk","time":250}`,
				}
				addMultipleActions(averager, actions)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"walk","avg":225}`,
					`{"action":"run","avg":125}`,
					`{"action":"crawl","avg":300}`,
				}
				verifyMultipleDifferentStats(stats, expStats)
			})

			It("should keep averaging after GetStats is called", func() {
				actions := []string{
					`{"action":"walk","time":225}`,
					`{"action":"run","time":75}`,
				}
				addMultipleActions(averager, actions)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"run","avg":75}`,
					`{"action":"walk","avg":225}`,
				}
				verifyMultipleDifferentStats(stats, expStats)
				err := averager.AddAction(`{"action":"run","time":80}`)
				Expect(err).NotTo(HaveOccurred())
				stats = averager.GetStats()
				expStats = []string{
					`{"action":"run","avg":77.5}`,
					expStats[1],
				}
				verifyMultipleDifferentStats(stats, expStats)
			})

			It("should handle an action with a time of 0", func() {
				err := averager.AddAction(`{"action":"bike","time":0}`)
				Expect(err).NotTo(HaveOccurred())
				stats := averager.GetStats()
				Expect(stats).To(Equal(`[{"action":"bike","avg":0}]`))
				err = averager.AddAction(`{"action":"bike","time":50}`)
				Expect(err).NotTo(HaveOccurred())
				stats = averager.GetStats()
				Expect(stats).To(Equal(`[{"action":"bike","avg":25}]`))
			})

			It("should not return anything if GetStats is called without AddAction being called", func() {
				stats := averager.GetStats()
				Expect(stats).To(Equal(emptyStats))
			})
		})

		Context("concurrent function calls", func() {
			It("should not average multiple inputs to different actions concurrently", func() {

			})

			It("should average multiple inputs of the same action concurrently", func() {

			})

			It("should average multiple inputs of different actions concurrently", func() {

			})

			It("should handle a mix of single and multiple inputs to different actions concurrently", func() {

			})

			It("should handle concurrent calls to AddAction and GetStats", func() {

			})
		})
	})

	Context("incorrect input provided", func() {
		It("should fail if an empty string is given as input", func() {
			err := averager.AddAction("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unexpected end of JSON input"))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if a non json string is given as input", func() {
			err := averager.AddAction("string")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid character"))
			Expect(err.Error()).To(ContainSubstring("looking for beginning of value"))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if values are not correct for the fields", func() {
			err := averager.AddAction(`{"action":123,"time":"run"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`action field is not string in input {"action":123,"time":"run"}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if fields are not correct", func() {
			err := averager.AddAction(`{"action":"run","tome":123}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`input {"action":"run","tome":123} is missing "time" field, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if there are additional fields", func() {
			err := averager.AddAction(`{"action":"run","time":123,"extra":true}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`extra field in input {"action":"run","time":123,"extra":true}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if there are duplicate fields", func() {
			err := averager.AddAction(`{"action":"run","action":"jump"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`extra field in input {"action":"run","action":"jump"}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if there is a negative value for time", func() {
			err := averager.AddAction(`{"action":"bike","time":-1}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`negative time value for input {"action":"bike","time":-1}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})
	})

	Context("incorrect and correct input", func() {
		It("should handle correct and incorrect calls", func() {
			actions := []string{
				`{"action":"run","time":20}`,
				`{"action":"jump","time":10}`,
			}
			addMultipleActions(averager, actions)
			err := averager.AddAction(`{"act":"jump","time":10}`)
			Expect(err).To(HaveOccurred())
			err = averager.AddAction(`{"action":"run","time":30}`)
			Expect(err).NotTo(HaveOccurred())
			stats := averager.GetStats()
			expStats := []string{
				`{"action":"run","avg":25}`,
				`{"action":"jump","avg":10}`,
			}
			verifyMultipleDifferentStats(stats, expStats)
		})
	})
})