package main

type encoding []byte

func (e encoding) Len() int           { return len(e) }
func (e encoding) Less(i, j int) bool { return i > j }
func (e encoding) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
