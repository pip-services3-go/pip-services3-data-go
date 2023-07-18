package test_persistence

type ItemPage struct {
	Total *int64 `json:"total"`
	Data  []Item `json:"data"`
}

func NewEmptyItemPage() *ItemPage {
	return &ItemPage{}
}

func NewItemPage(total *int64, data []Item) *ItemPage {
	return &ItemPage{Total: total, Data: data}
}
