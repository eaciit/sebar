package sebar

type IServer interface {
	Start() error
	Stop() error
}
