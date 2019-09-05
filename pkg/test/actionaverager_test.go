package actionaverager_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/action-averager/pkg/actionaverager"
)

const (
	emptyStats = "[]"

	delay     = true
	verifyAll = true

	replace1      = 1
	delayDuration = 10
)

func addMultipleActions(averager actionaverager.ActionAverager, actions []string, isDelay bool) {
	for i := range actions {
		err := averager.AddAction(actions[i])
		Expect(err).NotTo(HaveOccurred())
		// NOTE: delay should only be used when running concurrent tests
		if isDelay {
			time.Sleep(delayDuration * time.Millisecond)
		}
	}
}

func verifyMultipleDifferentStats(stats string, expSubstrings []string, fullVerify bool) {
	expNumCommas := len(expSubstrings) - 1
	for i := range expSubstrings {
		//NOTE: use ContainSubstring, since order of actions is go map order so explicit string equals fails intermittently
		Expect(stats).To(ContainSubstring(expSubstrings[i]))
		stats = strings.Replace(stats, expSubstrings[i], "", replace1)
	}
	// NOTE: this make sure the string only contains what is in the expected substrings and that it is properly formatted
	// meaning there should be 1 less comma than fields and only 1 [ and ]. This should not be run in some concurrent cases
	// where the exact output can not be guaranteed at any time i.e. concurrent calls of both AddAction and GetStats.
	if fullVerify {
		stats = strings.Replace(stats, ",", "", expNumCommas)
		stats = strings.Replace(stats, "[", "", replace1)
		stats = strings.Replace(stats, "]", "", replace1)
		Expect(stats).To(BeEmpty())
	}
}

var _ = Describe("action-averager tests", func() {
	// NOTE: create new averager for each test being run
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
				addMultipleActions(averager, actions, !delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"run","avg":55}`,
					`{"action":"skip","avg":145}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should average multiple inputs of the same action", func() {
				actions := []string{
					`{"action":"hop","time":55.5}`,
					`{"action":"hop","time":145.37}`,
				}
				addMultipleActions(averager, actions, !delay)
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
				addMultipleActions(averager, actions, !delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"hop","avg":59.025}`,
					`{"action":"skip","avg":145.38933333333333}`,
					`{"action":"jump","avg":32.785}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should handle a mix of single and multiple inputs to different actions", func() {
				actions := []string{
					`{"action":"walk","time":200}`,
					`{"action":"run","time":100}`,
					`{"action":"crawl","time":300}`,
					`{"action":"run","time":150}`,
					`{"action":"walk","time":250}`,
				}
				addMultipleActions(averager, actions, !delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"walk","avg":225}`,
					`{"action":"run","avg":125}`,
					`{"action":"crawl","avg":300}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should keep averaging after GetStats is called", func() {
				actions := []string{
					`{"action":"walk","time":225}`,
					`{"action":"run","time":75}`,
				}
				addMultipleActions(averager, actions, !delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"run","avg":75}`,
					`{"action":"walk","avg":225}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)

				err := averager.AddAction(`{"action":"run","time":80}`)
				Expect(err).NotTo(HaveOccurred())
				stats = averager.GetStats()
				expStats = []string{
					`{"action":"run","avg":77.5}`,
					expStats[1],
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
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
				actions0 := []string{
					`{"action":"bike","time":10}`,
					`{"action":"swim","time":20}`,
					`{"action":"run","time":30}`,
				}
				actions1 := []string{
					`{"action":"hop","time":40}`,
					`{"action":"skip","time":50}`,
					`{"action":"jump","time":60}`,
				}
				go addMultipleActions(averager, actions0, delay)
				addMultipleActions(averager, actions1, delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"bike","avg":10}`,
					`{"action":"swim","avg":20}`,
					`{"action":"run","avg":30}`,
					`{"action":"hop","avg":40}`,
					`{"action":"skip","avg":50}`,
					`{"action":"jump","avg":60}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should average multiple inputs of the same action concurrently", func() {
				actions0 := []string{
					`{"action":"bike","time":10}`,
					`{"action":"bike","time":20}`,
					`{"action":"bike","time":30}`,
				}
				actions1 := []string{
					`{"action":"bike","time":40}`,
					`{"action":"bike","time":50}`,
					`{"action":"bike","time":60}`,
				}
				go addMultipleActions(averager, actions1, delay)
				addMultipleActions(averager, actions0, delay)
				stats := averager.GetStats()
				Expect(stats).To(Equal(`[{"action":"bike","avg":35}]`))
			})

			It("should average multiple inputs of different actions concurrently", func() {
				actions0 := []string{
					`{"action":"bike","time":10}`,
					`{"action":"swim","time":20}`,
					`{"action":"bike","time":30}`,
				}
				actions1 := []string{
					`{"action":"swim","time":40}`,
					`{"action":"bike","time":50}`,
					`{"action":"swim","time":60}`,
				}
				go addMultipleActions(averager, actions0, delay)
				addMultipleActions(averager, actions1, delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"bike","avg":30}`,
					`{"action":"swim","avg":40}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should handle a mix of single and multiple inputs to different actions concurrently", func() {
				actions0 := []string{
					`{"action":"bike","time":10}`,
					`{"action":"swim","time":20}`,
					`{"action":"run","time":30}`,
				}
				actions1 := []string{
					`{"action":"swim","time":40}`,
					`{"action":"bike","time":50}`,
					`{"action":"walk","time":60}`,
				}
				go addMultipleActions(averager, actions0, delay)
				addMultipleActions(averager, actions1, delay)
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"bike","avg":30}`,
					`{"action":"swim","avg":30}`,
					`{"action":"run","avg":30}`,
					`{"action":"walk","avg":60}`,
				}
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
			})

			It("should handle concurrent calls to AddAction and GetStats", func() {
				// NOTE: verifications of the stats in these closures (funcs0 and funcs1) only verify
				// the minimum set of stats that can be expected, since this is run concurrently and
				// there is not any guarantee of what the stats actually look like when the verifications
				// are run. The minimum set chosen is from the set of actions that have occurred before
				// GetStats was called in the current closure, since these can be guaranteed.
				funcs0 := func() {
					actions0 := []string{
						`{"action":"bike","time":100}`,
						`{"action":"run","time":90}`,
					}
					actions1 := []string{
						`{"action":"swim","time":80}`,
						`{"action":"hop","time":70}`,
					}

					addMultipleActions(averager, actions0, delay)
					stats := averager.GetStats()
					minExpStats := []string{
						`{"action":"bike","avg":100}`,
						`{"action":"run","avg":90}`,
					}
					verifyMultipleDifferentStats(stats, minExpStats, !verifyAll)
					addMultipleActions(averager, actions1, delay)
				}

				funcs1 := func() {
					actions := []string{
						`{"action":"skip","time":60}`,
						`{"action":"jump","time":50}`,
					}
					err := averager.AddAction(`{"action":"walk","time":40}`)
					Expect(err).NotTo(HaveOccurred())
					stats := averager.GetStats()
					minExpStats0 := []string{`{"action":"walk","avg":40}`}
					verifyMultipleDifferentStats(stats, minExpStats0, !verifyAll)

					addMultipleActions(averager, actions, delay)
					stats = averager.GetStats()
					minExpStats1 := []string{
						`{"action":"walk","avg":40}`,
						`{"action":"skip","avg":60}`,
						`{"action":"jump","avg":50}`,
					}
					verifyMultipleDifferentStats(stats, minExpStats1, !verifyAll)
				}

				go funcs1()
				funcs0()
				stats := averager.GetStats()
				expStats := []string{
					`{"action":"bike","avg":100}`,
					`{"action":"run","avg":90}`,
					`{"action":"walk","avg":40}`,
					`{"action":"skip","avg":60}`,
					`{"action":"jump","avg":50}`,
					`{"action":"swim","avg":80}`,
					`{"action":"hop","avg":70}`,
				}
				// NOTE: full verification can happen here since all actions have stopped being added
				verifyMultipleDifferentStats(stats, expStats, verifyAll)
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
			Expect(err.Error()).To(Equal("invalid character 's' looking for beginning of value"))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if values are not correct for the fields", func() {
			err := averager.AddAction(`{"action":123,"time":"run"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`action field is not a string in input {"action":123,"time":"run"}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))

			err = averager.AddAction(`{"action":"run","time":"run"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`time field is not a number in input {"action":"run","time":"run"}, rejecting`))
			stats = averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if fields are not correct", func() {
			err := averager.AddAction(`{"action":"run","tome":123}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`input {"action":"run","tome":123} is missing "time" field, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))

			err = averager.AddAction(`{"actoin":"run","time":123}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`input {"actoin":"run","time":123} is missing "action" field, rejecting`))
			stats = averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if there are additional fields", func() {
			err := averager.AddAction(`{"action":"run","time":123,"extra":true}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`unexpected number of fields, 3, in input {"action":"run","time":123,"extra":true}, expect 2, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal(emptyStats))
		})

		It("should fail if there are duplicate fields", func() {
			err := averager.AddAction(`{"action":"run","action":"jump"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`unexpected number of fields, 1, in input {"action":"run","action":"jump"}, expect 2, rejecting`))
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
			addMultipleActions(averager, actions, !delay)
			err := averager.AddAction(`{"act":"jump","time":10}`)
			Expect(err).To(HaveOccurred())

			err = averager.AddAction(`{"action":"run","time":30}`)
			Expect(err).NotTo(HaveOccurred())
			stats := averager.GetStats()
			expStats := []string{
				`{"action":"run","avg":25}`,
				`{"action":"jump","avg":10}`,
			}
			verifyMultipleDifferentStats(stats, expStats, verifyAll)
		})
	})
})
