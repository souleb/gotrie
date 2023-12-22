package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Trie(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		checkFor func(*Trie) bool
	}{
		{
			name:  "test for an empty trie",
			input: []string{},
			checkFor: func(t *Trie) bool {
				return !t.HasItem("hello")
			},
		},
		{
			name:  "test for a single word",
			input: []string{"hello"},
			checkFor: func(t *Trie) bool {
				return t.HasItem("hello")
			},
		},
		{
			name:  "test for multiple words",
			input: []string{"hello", "world", "my-world"},
			checkFor: func(t *Trie) bool {
				return t.HasItem("hello") && t.HasItem("world") && t.HasItem("my-world")
			},
		},
		{
			name:  "test for multiple words with common prefix",
			input: []string{"home", "homework"},
			checkFor: func(t *Trie) bool {
				return t.HasItem("home") && t.HasItem("homework")
			},
		},
		{
			name:  "test for deleting a word",
			input: []string{"hello", "world", "home", "work", "homework"},
			checkFor: func(t *Trie) bool {
				t.DeleteItem("home")
				return t.HasItem("homework") && !t.HasItem("home")
			},
		},
		{
			name:  "test for getting a key's data",
			input: []string{"hello", "world", "home", "work", "homework"},
			checkFor: func(t *Trie) bool {
				return string(t.GetItem("homework")) == "homework-value" &&
					string(t.GetItem("hello")) == "hello-value" &&
					string(t.GetItem("world")) == "world-value" &&
					string(t.GetItem("home")) == "home-value" &&
					string(t.GetItem("work")) == "work-value"
			},
		},
		{
			name:  "test for setting a key's data",
			input: []string{"hello", "world", "home", "work", "homework"},
			checkFor: func(t *Trie) bool {
				t.Set("homework", []byte("homework grade 2"))
				return string(t.GetItem("homework")) == "homework grade 2"
			},
		},
		{
			name:  "test for getting all keys with a prefix",
			input: []string{"hello", "world", "home", "work", "homework"},
			checkFor: func(t *Trie) bool {
				keys := t.Keys("hom")
				return len(keys) == 2 && keys[0] == "home" && keys[1] == "homework"
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := NewTrie()
			for _, word := range tc.input {
				trie.InsertItem(word, []byte(word+"-value"))
			}
			assert.True(t, tc.checkFor(trie))
		})
	}
}
