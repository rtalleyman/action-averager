package actionaverager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("action-averager tests", func() {
	BeforeEach(func() {

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

		})

		It("should fail if a non json string is given as input", func() {

		})

		It("should fail if a not properly formatted json string is provided", func() {

		})
	})
})
