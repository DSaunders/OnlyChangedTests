package jest

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func WriteFilterFile(testsToRun []string) {
	log.Println(".. Generating test filter JavaScript file")
	err := ioutil.WriteFile("selected-tests.js", []byte(getTestFileContents(testsToRun)), 0644)
	if err != nil {
		panic(err)
	}
}

func RemoveFilterFile() {
	log.Println(".. Removing test filter JavaScript file")
	err := os.Remove("selected-tests.js")
	if err != nil {
		panic(err)
	}
}

func getTestFileContents(files []string) string {
	code := `
	const toRun = [
{{TESTS}}
	]

	module.exports = testPaths => {		
		return {
			filtered: testPaths.filter(a => {
				return toRun.includes(a.toLowerCase())
				})
				.map((testPath) => ({ test: testPath }))
			};
	};
`

	filesForJs := make([]string, len(files))
	for i, file := range files {
		filesForJs[i] = "        \"" + strings.Replace(file, "\\", "\\\\", -1) + "\","
	}

	code = strings.Replace(code, "{{TESTS}}", strings.Join(filesForJs, "\n"), 1)
	return code
}
