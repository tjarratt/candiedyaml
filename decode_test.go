package candiedyaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Decode", func() {
	It("Decodes a file", func() {
		f, _ := os.Open("fixtures/specification/example2_1.yaml")
		d := NewDecoder(f)
		v := new(interface{})
		err := d.Decode(&v)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("Sequence", func() {
		It("Decodes to interface{}s", func() {
			f, _ := os.Open("fixtures/specification/example2_1.yaml")
			d := NewDecoder(f)
			v := new(interface{})

			err := d.Decode(&v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω((*v).([]interface{})).To(Equal([]interface{}{"Mark McGwire", "Sammy Sosa", "Ken Griffey"}))
		})

		It("Decodes to []string", func() {
			f, _ := os.Open("fixtures/specification/example2_1.yaml")
			d := NewDecoder(f)
			v := make([]string, 0, 3)

			err := d.Decode(&v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(v).To(Equal([]string{"Mark McGwire", "Sammy Sosa", "Ken Griffey"}))
		})
	})

	Context("Mapping", func() {
		It("Decodes to interface{}s", func() {
			f, _ := os.Open("fixtures/specification/example2_2.yaml")
			d := NewDecoder(f)
			v := new(interface{})

			err := d.Decode(&v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω((*v).(map[interface{}]interface{})).To(Equal(map[interface{}]interface{}{
				"hr":  int64(65),
				"avg": float64(0.278),
				"rbi": int64(147),
			}))
		})

		It("Decodes to a map of string arrays", func() {
			f, _ := os.Open("fixtures/specification/example2_9.yaml")
			d := NewDecoder(f)
			v := make(map[string][]string)

			err := d.Decode(&v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(v).To(Equal(map[string][]string{}))
		})

	})

	Context("Sequence of Maps", func() {
		It("Decodes to interface{}s", func() {
			f, _ := os.Open("fixtures/specification/example2_4.yaml")
			d := NewDecoder(f)
			v := new(interface{})

			err := d.Decode(&v)
			Ω(err).ShouldNot(HaveOccurred())
			Ω((*v).([]interface{})).To(Equal([]interface{}{
				map[interface{}]interface{}{"name": "Mark McGwire", "hr": int64(65), "avg": float64(0.278)},
				map[interface{}]interface{}{"name": "Sammy Sosa", "hr": int64(63), "avg": float64(0.288)},
			}))
		})
	})

})
