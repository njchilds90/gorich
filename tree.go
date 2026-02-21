package gorich

import (
	"fmt"
	"io"
	"os"
)

// TreeNode represents a node in a tree.
type TreeNode struct {
	label    string
	style    Style
	children []*TreeNode
}

// NewTreeNode creates a new TreeNode.
func NewTreeNode(label string, style ...Style) *TreeNode {
	n := &TreeNode{label: label}
	if len(style) > 0 {
		n.style = style[0]
	}
	return n
}

// Add appends a child node and returns the child (for chaining).
func (n *TreeNode) Add(label string, style ...Style) *TreeNode {
	child := NewTreeNode(label, style...)
	n.children = append(n.children, child)
	return child
}

// AddNode appends an existing node as a child.
func (n *TreeNode) AddNode(child *TreeNode) *TreeNode {
	n.children = append(n.children, child)
	return child
}

// Tree renders a tree structure to the terminal.
type Tree struct {
	root        *TreeNode
	guideStyle  Style
	labelStyle  Style
	w           io.Writer
}

// NewTree creates a tree with a root label.
func NewTree(rootLabel string, opts ...func(*Tree)) *Tree {
	t := &Tree{
		root:       NewTreeNode(rootLabel),
		guideStyle: NewStyle(Dim, BrightBlue),
		labelStyle: NewStyle(Bold),
		w:          os.Stdout,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Tree option setters.

func TreeGuideStyle(s Style) func(*Tree)  { return func(t *Tree) { t.guideStyle = s } }
func TreeLabelStyle(s Style) func(*Tree)  { return func(t *Tree) { t.labelStyle = s } }
func TreeWriter(w io.Writer) func(*Tree)  { return func(t *Tree) { t.w = w } }

// Root returns the root node for adding children.
func (t *Tree) Root() *TreeNode { return t.root }

// Render prints the tree.
func (t *Tree) Render() {
	label := t.root.label
	if t.root.style.codes != nil {
		label = t.root.style.Apply(label)
	} else {
		label = t.labelStyle.Apply(label)
	}
	fmt.Fprintln(t.w, label)
	t.renderChildren(t.root.children, "")
}

func (t *Tree) renderChildren(children []*TreeNode, prefix string) {
	for i, child := range children {
		isLast := i == len(children)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		label := child.label
		if child.style.codes != nil {
			label = child.style.Apply(label)
		}
		fmt.Fprintf(t.w, "%s%s%s\n",
			t.guideStyle.Apply(prefix+connector),
			label,
			"",
		)
		t.renderChildren(child.children, childPrefix)
	}
}

// PrintTree is a convenience wrapper.
func PrintTree(rootLabel string, build func(*TreeNode), opts ...func(*Tree)) {
	tr := NewTree(rootLabel, opts...)
	build(tr.Root())
	tr.Render()
}
