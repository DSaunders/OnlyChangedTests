package dependencytree

type Node struct {
	FileName string
	Parents  []*Node
	Children []*Node
}

func MakeNode(filename string) *Node {
	return &Node{
		FileName: filename,
		Parents:  make([]*Node, 0),
		Children: make([]*Node, 0),
	}
}
