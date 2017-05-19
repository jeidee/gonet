package gonet

type Stub interface {
	ParseAndDo(data *IncommingData) error
}
