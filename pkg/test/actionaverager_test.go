package actionaverager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/action-averager/pkg/actionaverager"
)

var _ = Describe("action-averager tests", func() {
	var averager actionaverager.ActionAverager
	BeforeEach(func() {
		averager = actionaverager.NewActionAverager()
	})

	Context("correct input provided", func() {
		Context("serial function calls", func() {
			It("should not average a single input", func() {

			})

			It("should not average multiple inputs to different actions", func() {

			})

			It("should average multiple inputs of the same action", func() {

			})

			It("should average multiple inputs of different actions", func() {

			})

			It("should handle a mix of single and multiple inputs to different actions", func() {

			})

			It("should keep averaging after GetStats is called", func() {

			})

			It("should not return anything if GetStats is called without AddAction being called", func() {

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
			Expect(stats).To(Equal("[]"))
		})

		It("should fail if a non json string is given as input", func() {
			err := averager.AddAction("string")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid character"))
			Expect(err.Error()).To(ContainSubstring("looking for beginning of value"))
			stats := averager.GetStats()
			Expect(stats).To(Equal("[]"))
		})

		It("should fail if values are not correct for the fields", func() {
			err := averager.AddAction(`{"action":123,"time":"run"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`action field is not string in input {"action":123,"time":"run"}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal("[]"))
		})

		It("should fail if fields are not correct", func() {
			err := averager.AddAction(`{"action":"run","tome":123}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`input {"action":"run","tome":123} is missing "time" field, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal("[]"))
		})

		It("should fail if there are additional fields", func() {
			err := averager.AddAction(`{"action":"run","time":123,"extra":true}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`extra field in input {"action":"run","time":123,"extra":true}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal("[]"))
		})

		It("should fail if there are duplicate fields", func() {
			err := averager.AddAction(`{"action":"run","action":"jump"}`)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`extra field in input {"action":"run","action":"jump"}, rejecting`))
			stats := averager.GetStats()
			Expect(stats).To(Equal("[]"))
		})
	})

	Context("incorrect and correct input", func() {
		It("should handle correct and incorrect calls", func() {
			err := averager.AddAction(`{"action":"run","time":20}`)
			Expect(err).NotTo(HaveOccurred())
			err = averager.AddAction(`{"action":"jump","time":10}`)
			Expect(err).NotTo(HaveOccurred())
			err = averager.AddAction(`{"act":"jump","time":10}`)
			Expect(err).To(HaveOccurred())
			err = averager.AddAction(`{"action":"run","time":30}`)
			Expect(err).NotTo(HaveOccurred())
			stats := averager.GetStats()
			Expect(stats).To(Equal(`[{"action":"run","avg":25},{"action":"jump","avg":10}]`))
		})
	})
})
