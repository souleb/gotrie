//This package implements a trie data structure as proposed in https://en.wikipedia.org/wiki/Radix_tree.
//The trie is not thread-safe.

package trie

import (
	"sort"
)

// TFunc is the type of the function called for each node.
type TFunc func(step int, path string, data []byte, isLeaf bool) (bool, error)

// Trie is a trie data structure.
type Trie struct {
	root   *node
	sorted bool
}

type node struct {
	isLeaf bool
	edges  map[byte]*edge
	data   []byte
}

type edge struct {
	label []byte
	next  *node
}

// NewTrie returns a new Trie.
func NewTrie() *Trie {
	return &Trie{
		root: &node{
			false,
			make(map[byte]*edge),
			nil,
		},
	}
}

// Sorted sets the sorted flag.
// If sorted is true, the edges will be sorted by their keys during traversal.
func (t *Trie) Sorted(sorted bool) {
	t.sorted = sorted
}

// InsertItem inserts the given item into the trie.
func (t *Trie) InsertItem(key string, data []byte) {
	cur := t.root
	bPref := []byte(key)
	for len(bPref) > 0 {
		k := bPref[0]
		if cur.edges == nil {
			cur.edges = make(map[byte]*edge)
		}
		currEdge, exists := cur.edges[k]
		if !exists {
			cur.edges[k] = &edge{
				label: bPref,
				next: &node{
					isLeaf: true,
					data:   data,
				},
			}
			break
		}
		curStr, curStrLen := bPref, len(bPref)
		labelLen := len(currEdge.label)
		// if the current string is longer than the edge label
		// then we split it to fit the edge label for this pass
		if curStrLen > labelLen {
			curStr = bPref[:labelLen]
		}

		splitIdx := getFirstMismatch(curStr, currEdge.label)
		if splitIdx != labelLen {
			tail := currEdge.label[splitIdx:]
			currEdge.label = currEdge.label[:splitIdx]
			newNext := &node{false, make(map[byte]*edge), nil}
			newNext.edges[tail[0]] = &edge{
				label: tail,
				next:  currEdge.next,
			}
			currEdge.next = newNext
		}

		if len(bPref) == len(currEdge.label) {
			currEdge.next.isLeaf = true
			currEdge.next.data = data
		}

		cur = currEdge.next
		bPref = bPref[splitIdx:]
	}
}

// getFirstMismatch returns the index of the first mismatch between current and
// label.
func getFirstMismatch(current, label []byte) int {
	if len(current) > len(label) {
		current, label = label, current
	}
	// iterate over the shorter string
	for i := 0; i < len(current); i++ {
		if current[i] != label[i] {
			return i
		}
	}
	return len(current)
}

// getNode returns the node for the given prefix.
func (t *Trie) getNode(prefix string) *node {
	cur := t.root
	bPref := []byte(prefix)
	for len(bPref) > 0 {
		edge, exists := cur.edges[bPref[0]]
		if !exists {
			return nil
		}
		splitIdx := getFirstMismatch(bPref, edge.label)
		if splitIdx != len(edge.label) {
			return nil
		}
		cur = edge.next
		bPref = bPref[splitIdx:]
	}
	return cur
}

// HasItem returns true if the given key exists in the trie.
func (t *Trie) HasItem(key string) bool {
	node := t.getNode(key)
	return node != nil && node.isLeaf
}

// DeleteItem deletes the item with the given key.
func (t *Trie) DeleteItem(key string) {
	t.root = t.delete(t.root, []byte(key))
}

func (t *Trie) delete(node *node, key []byte) *node {
	if len(key) == 0 {
		if node.edges == nil && node != t.root {
			return nil
		}
		node.isLeaf = false
		return node
	}

	currEdge, exists := node.edges[key[0]]
	if !exists {
		return node
	}

	deleted := t.delete(currEdge.next, key[len(currEdge.label):])
	if deleted == nil {
		delete(node.edges, key[0])
		if len(node.edges) == 0 && !node.isLeaf && node != t.root {
			return nil
		}
	} else if len(deleted.edges) == 1 && !deleted.isLeaf {
		delete(node.edges, key[0])
		for _, v := range deleted.edges {
			node.edges[key[0]] = &edge{
				label: append(currEdge.label, v.label...),
				next:  v.next,
			}
		}
	}

	return node
}

// GetItem returns the data for the given item.
func (t *Trie) GetItem(key string) []byte {
	return t.getNode(key).data
}

// SetI sets the data for the given item.
func (t *Trie) Set(key string, data []byte) {
	t.getNode(key).data = data
}

// Keys returns all the keys in the trie that start with the given prefix.
func (t *Trie) Keys(prefix string) []string {
	node, pref := t.startsWith(prefix)
	if node == nil {
		return nil
	}
	var keys []string
	if node.isLeaf {
		keys = append(keys, pref)
	}
	for _, edge := range node.edges {
		keys = append(keys, t.Keys(pref+string(edge.label))...)
	}
	return keys
}

// startsWith returns the node for the prefix that starts with the given prefix.
func (t *Trie) startsWith(prefix string) (*node, string) {
	cur := t.root
	bPref := []byte(prefix)
	for len(bPref) > 0 {
		edge, exists := cur.edges[bPref[0]]
		if !exists {
			return nil, ""
		}
		splitIdx := getFirstMismatch(bPref, edge.label)
		if splitIdx != len(edge.label) {
			if splitIdx == len(bPref) {
				return edge.next, prefix + string(edge.label[splitIdx:])
			}
			return nil, ""
		}
		cur = edge.next
		bPref = bPref[splitIdx:]
	}
	return cur, prefix
}

// Traverse traverses the trie and calls the given function for each node.
func (t *Trie) Traverse(fn TFunc) error {
	return t.traverse(t.root, 0, []byte{}, fn)
}

// traverse recursively traverses the trie and calls the given function for each node.
func (t *Trie) traverse(node *node, step int, path []byte, fn TFunc) error {
	if node == nil {
		return nil
	}

	next, err := fn(step, string(path), node.data, node.isLeaf)
	if err != nil {
		return err
	}

	// it's time to stop
	if !next {
		return nil
	}

	keys := keys(node.edges, t.sorted)
	for _, k := range keys {
		err = t.traverse(node.edges[k].next, step+1, append(path, node.edges[k].label...), fn)
		if err != nil {
			return err
		}
	}
	return nil
}

// keys returns the keys of the given map.
// If sorted is true, the keys will be sorted.
func keys(m map[byte]*edge, sorted bool) []byte {
	var (
		keys []byte
	)
	for k := range m {
		keys = append(keys, k)
	}

	if sorted {
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	}
	return keys
}
