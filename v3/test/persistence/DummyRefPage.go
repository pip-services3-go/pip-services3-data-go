package test_persistence

type DummyRefPage struct {
	Total *int64   `json:"total"`
	Data  []*Dummy `json:"data"`
}

func NewEmptyDummyRefPage() *DummyRefPage {
	return &DummyRefPage{}
}

func NewDummyRefPage(total *int64, data []*Dummy) *DummyRefPage {
	return &DummyRefPage{Total: total, Data: data}
}
