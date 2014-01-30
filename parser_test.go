package candiedyaml

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
)

var parses = func(filename string) {
	It("parses "+filename, func() {
		file, err := os.Open(filename)
		Ω(err).To(BeNil())

		parser := yaml_parser_t{}
		yaml_parser_initialize(&parser)
		yaml_parser_set_input_file(&parser, file)

		failed := false
		token := yaml_token_t{}

		for {
			if !yaml_parser_scan(&parser, &token) {
				failed = true
				break
			}

			if token.token_type == YAML_STREAM_END_TOKEN {
				break
			}
		}

		file.Close()

		// msg := "SUCCESS"
		// if failed {
		// 	msg = "FAILED"
		// 	if parser.error != YAML_NO_ERROR {
		// 		m := parser.problem_mark
		// 		fmt.Printf("ERROR: (%s) %s @ line: %d  col: %d\n",
		// 			parser.context, parser.problem, m.line, m.column)
		// 	}
		// }
		Ω(failed).To(BeFalse())
	})
}

var testYaml = func(dirname string) {
	fileInfos, err := ioutil.ReadDir(dirname)
	Ω(err).To(BeNil())
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			parses(filepath.Join(dirname, fileInfo.Name()))
		}
	}
}

var _ = Describe("Parser", func() {
	testYaml("fixtures/specification")
	testYaml("fixtures/specification/types")
})
