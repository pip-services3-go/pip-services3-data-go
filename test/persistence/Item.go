package test_persistence

type Item struct {
	Id                        string `json:"id"`
	FailingToUpdateThisField1 int64  `json:"failing_to_update_this_field_1,string"`
	FailingToUpdateThisField2 int64  `json:"failing_to_update_this_field_2,string"`
	UpdatedBy                 string `json:"updated_by"`
}
