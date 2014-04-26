package main

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

type Shortcut struct {
	Id  int64
	URL string
}

var (
	alphabet  = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	shortcuts = make(map[int64]*Shortcut)
)

var ErrNotFound = errors.New("shortcut not found")

func NewShortcut() *Shortcut {
	return &Shortcut{}
}

func FindShortcut(encoded string) (*Shortcut, error) {
	id := 0
	base := len(alphabet)
	for _, v := range encoded {
		id = id*base + bytes.IndexByte(alphabet, byte(v))
	}

	shortcut, ok := shortcuts[int64(id)]
	if !ok {
		return nil, ErrNotFound
	}

	return shortcut, nil
}

func (s *Shortcut) Encode() string {
	if s.Id == 0 {
		return ""
	}

	rv := encoding{}
	base := int64(len(alphabet))

	for i := s.Id; i > 0; i /= base {
		rv = append(rv, alphabet[i%base])
	}

	sort.Sort(rv)
	return string(rv)
}

func (s *Shortcut) Save() {
	if s.Id == 0 {
		s.Id = int64(len(shortcuts)) + 1
	}

	shortcuts[s.Id] = s
}

func (s *Shortcut) String() string {
	return fmt.Sprintf("<Shortcut Id:%d, URL:%q", s.Id, s.URL)
}
