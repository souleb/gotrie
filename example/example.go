// This example shows how to use the trie package to store file paths.
//
// The file paths are stored as keys and the file type (file or directory) as data.
//
// The PrintFiles method prints all the file paths in the trie.
// The output of this example is:
//   - /etc/nginx/nginx.conf
//   - /etc/sshd/sshd.conf
//   - /var/log/nginx/access.log

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	trie "github.com/souleb/gotrie"
)

const (
	dir  = "dir"
	file = "file"
)

// StringTrie is a wrapper around trie.Trie.
type StringTrie struct {
	trie.Trie
	out io.Writer
}

// NewStringTrie returns a new StringTrie.
func NewStringTrie(out io.Writer) *StringTrie {
	return &StringTrie{
		Trie: *trie.NewTrie(),
		out:  out,
	}
}

// PrintFiles prints all the files in the trie.
func (t *StringTrie) PrintFiles() error {
	err := t.Traverse(t.printFile)
	if err != nil {
		return err
	}
	return nil
}

// printFile prints the file path if the node is a leaf and its data is "file".
// it implements trie.TFunc.
func (t *StringTrie) printFile(step int, path string, data []byte, isLeaf bool) (bool, error) {
	if bytes.Equal(data, []byte(file)) && isLeaf {
		fmt.Fprintln(t.out, path)
	}
	return true, nil
}

func main() {
	t := NewStringTrie(os.Stdout)
	t.InsertItem("/etc/nginx/nginx.conf", []byte(file))
	t.InsertItem("/etc/nginx/nginx.conf.d", []byte(dir))
	t.InsertItem("/etc/sshd", []byte(dir))
	t.InsertItem("/etc/sshd/sshd.conf", []byte(file))
	t.InsertItem("/var/log", []byte(dir))
	t.InsertItem("/var/log/nginx", []byte(dir))
	t.InsertItem("/var/log/nginx/access.log", []byte(file))
	t.PrintFiles()
}
