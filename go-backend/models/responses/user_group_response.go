package responses

type UserGroupResponse struct {
	GroupId   uint   `json:"id" gorm:"column:id"`
	GroupName string `json:"group_name"`
	ImageUrl  string `json:"image_url"`
}
