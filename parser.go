package main

type DictEntry struct {
	Word       string
	Definition []string
}

type DictionaryParser interface {
	Entries(filepaths []string) <-chan DictEntry
}
