package dependencytree

import (
	"io/ioutil"
	"log"
	"onlychangedtests/config"
	"onlychangedtests/filelist"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type results struct {
	ImpactedTests []string
}

type Tree struct {
	testNodes   []*Node
	allNodes    map[string]*Node
	importRegex *regexp.Regexp
	fileList    *filelist.FileList
	config      *config.Config
}

func BuildForFiles(tests []string, filelist *filelist.FileList, config *config.Config) *Tree {
	newTree := Tree{
		fileList:    filelist,
		config:      config,
		importRegex: regexp.MustCompile(`(?im)import\s+?(?:(?:(?:[\w*\s{},]*)\s+from\s+?)|)(?:['|"](.*?)['|"]|(?:'.*?'))[\s]*?(?:;|$|)`),
	}

	newTree.Build(tests)

	return &newTree
}

func (tree *Tree) Build(rootPaths []string) {

	visited := make(map[string]*Node)
	graph := MakeNode("(root)")

	for _, testFile := range rootPaths {
		tree.recursivelyGetImports(visited, graph, testFile)
	}

	for _, node := range graph.Children {
		node.Parents = nil
	}

	tree.testNodes = graph.Children
	tree.allNodes = visited
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

		withFileType := tree.getFullFilenameForImport(filename, match)
		if withFileType == "" {
			continue
		}

		tree.recursivelyGetImports(visited, thisNode, withFileType)
	}
}

func (tree *Tree) getFullFilenameForImport(sourceFile string, importLocation string) string {
	validExensions := tree.config.ModuleFileExtensions

	var fullFilename string

	// Can only handle these types of imports for now
	if !strings.HasPrefix(importLocation, "./") && !strings.HasPrefix(importLocation, "../") {
		return fullFilename
	}

	sourceFileDirectory := filepath.Dir(sourceFile)
	relativePath := filepath.Join(sourceFileDirectory, importLocation)

	fileExtension := filepath.Ext(relativePath)

	if fileExtension != "" &&
		slices.Contains(validExensions, strings.ToLower(fileExtension)) {
		// File has a valid extension aready, try to find it as-is before
		// appending other extensions to it below
		if !tree.fileList.Exists(relativePath) {
			return relativePath
		}
	}

	// Try to add valid extensions to this import until
	// we find a real file
	for _, filetype := range validExensions {
		newName := relativePath + filetype

		if !tree.fileList.Exists(newName) {
			continue
		}

		fullFilename = newName
	}

	return fullFilename
}

func (tree *Tree) getImports(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	all := tree.importRegex.FindAllStringSubmatch(string(content), -1)

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

	if _, found := tree.allNodes[filename]; !found {
		// file isn't in our map, no test dependencies
		return make([]string, 0)
	}

	// Walk up all parents
	visited := make(map[*Node]bool, 0)

	results := results{
		ImpactedTests: make([]string, 0),
	}
	tree.visitNode(tree.allNodes[filename], &results, visited)

	return results.ImpactedTests
}

func (tree *Tree) visitNode(node *Node, tests *results, visited map[*Node]bool) {

	if _, found := visited[node]; found {
		// Already been here, do nothing
		return
	}

	visited[node] = true

	if len(node.Parents) == 0 {
		// I am the top level, so I am a test that must be run
		tests.ImpactedTests = append(tests.ImpactedTests, node.FileName)
	} else {
		// Not top level, recursively walk up the tree
		for _, parent := range node.Parents {
			tree.visitNode(parent, tests, visited)
		}
	}
}
