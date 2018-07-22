package util

import "github.com/mitchellh/hashstructure"

func CompareStruct(v1, v2 interface{}) bool {
	hash1, _ := hashstructure.Hash(v1, nil)
	hash2, _ :=hashstructure.Hash(v2, nil)
	return  hash1 == hash2
}
