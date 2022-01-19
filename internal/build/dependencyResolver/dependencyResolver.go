package dependencyResolver

import (
	"errors"
	"fmt"
)

type node struct {
	key      string
	exists   bool
	parent   *node
	children []*node
}

var dependencyTree *node
var dependencyIndex map[string]*node

func (n *node) GetKey() string {
	return n.key
}

func findKeyMatch(n *node, key string) bool {
	fmt.Printf("Checking if %v == %v\n", n.key, key)
	if n.key == key {
		return true
	}
	match := false
	for i := 0; i < len(n.children) && !match; i++ {
		match = findKeyMatch(n.children[i], key)
	}
	return false
}

func deleteChild(parent *node, child *node) {
	match := false
	var i int
	for i = 0; i < len(parent.children) && !match; i++ {
		if parent.children[i] == child {
			match = true
		}
	}

	if match {
		i--
		if i < len(parent.children)-1 {
			parent.children = append(parent.children[:i], parent.children[i+1:]...)
		} else {
			parent.children = parent.children[:i]
		}
	}
}

func Initialize() {
	dependencyTree = &node{}
	dependencyTree.key = "root"
	dependencyTree.exists = true

	dependencyIndex = make(map[string]*node)
}

func AddKey(key string) {
	childNode, exists := dependencyIndex[key]
	if !exists {
		childNode = &node{}
		childNode.key = key
		childNode.exists = true
		dependencyIndex[key] = childNode

		childNode.parent = dependencyTree
		dependencyTree.children = append(dependencyTree.children, childNode)
	}
}

func keyExists(key string) bool {
	_, exists := dependencyIndex[key]
	return exists
}

func getNode(key string) *node {
	return dependencyIndex[key]
}

func AddDependsOn(key string, dependsOn string) error {
	if !keyExists(key) {
		AddKey(key)
	}
	var childNode = getNode(key)
	childNode.exists = true

	n, exists := dependencyIndex[dependsOn]
	if !exists {
		n = &node{}
		n.key = dependsOn
		n.parent = dependencyTree
		n.exists = false
		dependencyIndex[dependsOn] = n

		dependencyTree.children = append(dependencyTree.children, n)
	}

	match := false

	// check if childNode has node as a descendent
	for i := 0; i < len(childNode.children) && !match; i++ {
		match = findKeyMatch(childNode.children[i], dependsOn)
	}
	if match {
		return errors.New("dependency loop found")
	} else {
		// remove childNode from old parent
		if childNode.parent != nil && childNode.parent != n {
			deleteChild(childNode.parent, childNode)
		}
		childNode.parent = n
		n.children = append(n.children, childNode)
	}

	return nil
}

func SetNoDepends(key string) {
	if !keyExists(key) {
		AddKey(key)
	}
	var childNode = getNode(key)
	childNode.exists = true
}

func Validate() error {
	return validateChildren(dependencyTree)
}

func validateChildren(n *node) error {
	var err error = nil

	if !n.exists {
		err = errors.New("Missing dependency " + n.key)
	} else {
		var i int
		for i = 0; i < len(n.children) && err == nil; i++ {
			err = validateChildren(n.children[i])
		}
	}
	return err
}

func PrintTree() {
	printChildren(dependencyTree, 0)
}

func printChildren(n *node, depth int) {
	var i int
	for i = 0; i < depth; i++ {
		fmt.Print(" ")
	}
	fmt.Println(n.key)
	for i = 0; i < len(n.children); i++ {
		printChildren(n.children[i], depth+1)
	}
}
