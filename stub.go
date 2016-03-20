package gonet

type Stub interface {
	ParseAndDo(data *IncommingData, err error)
}
