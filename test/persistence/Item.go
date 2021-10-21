package test_persistence

type Item struct {
	Id                             string `json:"id"`
	Failing_to_update_this_field_1 int64  `json:"failing_to_update_this_field_1,string"`
	Failing_to_update_this_field_2 int64  `json:"failing_to_update_this_field_2,string"`
	Updated_by                     string `json:"updated_by"`
}
