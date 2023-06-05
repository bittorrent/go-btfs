package corehttp

// import (
// 	"fmt"
// 	"strings"

// 	files "github.com/bittorrent/go-btfs-files"
// 	gopath "path"
// )

// func findNode(root files.Node, path string) (files.Node, error) {
// 	// Error/Edge cases
// 	path = gopath.Clean(path)
// 	if path == "" {
// 		return nil, fmt.Errorf("The given path string is empty")
// 	}

// 	// Define vars
// 	rootDir, ok := root.(files.Directory)
// 	if !ok {
// 		return nil, fmt.Errorf("The given root is not a directory")
// 	}
// 	parts := strings.Split(path, "/")
// 	i := 3

// 	// Find the node with part[i] as name and, if not the last, go down.
// 	it := rootDir.Entries()
// 	it.BreadthFirstTraversal()
// 	return visit(it, parts, i, path)
// }

// func findHelper(node files.Node, parts []string, i int, path string) (files.Node, error) {
// 	// Define vars
// 	dir, ok := node.(files.Directory)
// 	if !ok {
// 		return nil, fmt.Errorf("The given node is not a directory")
// 	}

// 	// Find the current part's node and go down.
// 	it := dir.Entries()
// 	it.BreadthFirstTraversal()
// 	return visit(it, parts, i, path)
// }

// func visit(it files.DirIterator, parts []string, i int, path string) (files.Node, error) {
// 	// Find the current part's node and go down.
// 	for it.Next() {
// 		if it.Name() == parts[i] {
// 			if len(parts) == i+1 {
// 				return it.Node(), nil
// 			}
// 			i++
// 			return findHelper(it.Node(), parts, i, path)
// 		}
// 	}
// 	return nil, fmt.Errorf("%s not found from the directory object: %s", path, it.Err())
// }
