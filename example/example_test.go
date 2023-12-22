package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StringTrie(t *testing.T) {
	var buf bytes.Buffer
	type ft struct {
		path string
		typ  string
	}
	testCases := []struct {
		name     string
		input    []ft
		checkFor func(*testing.T, *StringTrie) bool
	}{
		{
			name:  "test for a single directory",
			input: []ft{{"/tmp", dir}},
			checkFor: func(_ *testing.T, t *StringTrie) bool {
				return t.HasItem("/tmp") && string(t.GetItem("/tmp")) == dir
			},
		},
		{
			name:  "test for multiple directories",
			input: []ft{{"/tmp", dir}, {"/etc", dir}, {"/etc/systemd", dir}},
			checkFor: func(_ *testing.T, t *StringTrie) bool {
				return t.HasItem("/tmp") && string(t.GetItem("/tmp")) == dir &&
					t.HasItem("/etc") && string(t.GetItem("/etc")) == dir &&
					t.HasItem("/etc/systemd") && string(t.GetItem("/etc/systemd")) == dir
			},
		},
		{
			name:  "test for file and directory",
			input: []ft{{"/etc/nginx/nginx.conf", file}, {"/etc/nginx/nginx.conf.d", dir}},
			checkFor: func(_ *testing.T, t *StringTrie) bool {
				return t.HasItem("/etc/nginx/nginx.conf") && string(t.GetItem("/etc/nginx/nginx.conf")) == file &&
					t.HasItem("/etc/nginx/nginx.conf.d") && string(t.GetItem("/etc/nginx/nginx.conf.d")) == dir
			},
		},
		{
			name: "test for printing trie files",
			input: []ft{
				{"/etc/nginx/nginx.conf", file},
				{"/etc/nginx/nginx.conf.d", dir},
				{"/etc/sshd", dir},
				{"/etc/sshd/sshd.conf", file},
				{"/var/log", dir},
				{"/var/log/nginx", dir},
				{"/var/log/nginx/access.log", file},
			},
			checkFor: func(t *testing.T, trie *StringTrie) bool {
				require.NoError(t, trie.PrintFiles())
				assert.Equal(t, buf.String(), "/etc/nginx/nginx.conf\n"+
					"/etc/sshd/sshd.conf\n"+
					"/var/log/nginx/access.log\n")
				return true
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			trie := NewStringTrie(&buf)
			for _, v := range tc.input {
				trie.InsertItem(v.path, []byte(v.typ))
			}
			assert.True(t, tc.checkFor(t, trie))
		})
	}
}
