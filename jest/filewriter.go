package jest

import (
	"io/ioutil"
	"strings"
)

func WriteFilterFile(testsToRun []string) {
	err := ioutil.WriteFile("selected-tests.js", []byte(getTestFileContents(testsToRun)), 0644)
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
