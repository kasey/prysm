package sszgen

type SSZSatisfier interface {
	Imports() map[string]string
	Methods() map[string]string
}