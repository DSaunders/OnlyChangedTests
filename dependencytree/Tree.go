package dependencytree

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type results struct {
	ImpactedTests []string
}

type Tree struct {
	TestNodes   []*Node
	AllNodes    map[string]*Node
	ImportRegex *regexp.Regexp
}

func BuildForFiles(tests []string) *Tree {
	newTree := Tree{
		ImportRegex: regexp.MustCompile(`(?im)import\s+?(?:(?:(?:[\w*\s{},]*)\s+from\s+?)|)(?:['|"](.*?)['|"]|(?:'.*?'))[\s]*?(?:;|$|)`),
	}

	newTree.Build(tests)

	return &newTree
}

// TODO: Return the file tree
func (tree *Tree) Build(rootPaths []string) {

	visited := make(map[string]*Node)
	graph := MakeNode("(root)")

	for _, testFile := range rootPaths {
		tree.recursivelyGetImports(visited, graph, testFile)
	}

	for _, node := range graph.Children {
		node.Parents = nil
	}

	tree.TestNodes = graph.Children
	tree.AllNodes = visited
}

func (tree *Tree) recursivelyGetImports(visited map[string]*Node, parentNode *Node, filename string) {

	// Have we visited this file already?
	if node, found := visited[filename]; found {
		// We have, so add it to the parent
		node.Parents = append(node.Parents, parentNode)
		parentNode.Children = append(parentNode.Children, node)

		// We've already visited its children, we can stop here
		return
	}

	thisNode := MakeNode(filename)
	thisNode.Parents = append(thisNode.Parents, parentNode)
	parentNode.Children = append(parentNode.Children, thisNode)

	visited[filename] = thisNode

	imports := tree.getImports(filename)
	for _, match := range imports {

		withFileType := tree.getValidFile(filename, match)
		if withFileType == "" {
			continue
		}

		tree.recursivelyGetImports(visited, thisNode, withFileType)
	}

}

// optimisation, we've read the whole file tree ahead of time
// anyway to find the tests, so we can do this bit in memory
func (tree *Tree) getValidFile(rootFileName string, relativeFile string) string {
	var withFileType string

	// Is it an absolute import? If so, do nothing (can we resolve these somehow?)
	if !strings.HasPrefix(relativeFile, "./") && !strings.HasPrefix(relativeFile, "../") {
		return withFileType
	}

	dir := filepath.Dir(rootFileName)
	relative := filepath.Join(dir, relativeFile)

	// TODO: get from config
	fileTypes := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, filetype := range fileTypes {

		// TODO: if no extension already of course!
		newName := relative + filetype

		_, err := os.OpenFile(newName, os.O_RDONLY, 0400)
		if err != nil {
			continue
		}

		withFileType = newName
	}

	return withFileType
}

// TODO: Doesn't support commonjs yet
func (tree *Tree) getImports(filename string) []string {
	// Potential optimisation (test it), stream the file instead
	// 'imports' are always at the start, so if we find them we
	// can stop after them (in common js require can be anywhere though)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	all := tree.ImportRegex.FindAllStringSubmatch(string(content), -1)

	matches := make([]string, len(all))
	for i, match := range all {
		matches[i] = match[1]
	}

	return matches
}

func (tree *Tree) GetTopLevelNodesForFiles(filenames []string) []string {
	disinctTopLevelFiles := make([]string, 0)

	for _, file := range filenames {
		topLevelFiles := tree.GetTopLevelNodesForFile(file)

		for _, topLevelFile := range topLevelFiles {
			if !slices.Contains(disinctTopLevelFiles, strings.ToLower(topLevelFile)) {
				disinctTopLevelFiles = append(disinctTopLevelFiles, strings.ToLower(topLevelFile))
			}
		}
	}

	return disinctTopLevelFiles
}

func (tree *Tree) GetTopLevelNodesForFile(filename string) []string {

	if _, found := tree.AllNodes[filename]; !found {
		// file isn't in our map, no test dependencies
		return make([]string, 0)
	}

	// Walk up all parents
	visited := make(map[*Node]bool, 0)

	results := results{
		ImpactedTests: make([]string, 0),
	}
	tree.visitNode(tree.AllNodes[filename], &results, visited)

	return results.ImpactedTests
}

func (tree *Tree) visitNode(node *Node, tests *results, visited map[*Node]bool) {

	if _, found := visited[node]; found {
		// Already been here, do nothing
		return
	}

	visited[node] = true

	if len(node.Parents) == 0 {
		// I am the top level
		tests.ImpactedTests = append(tests.ImpactedTests, node.FileName)
	} else {
		for _, parent := range node.Parents {
			tree.visitNode(parent, tests, visited)
		}
	}
}
