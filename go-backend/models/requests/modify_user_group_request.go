package requests

type ModifyUserGroupRequest struct {
	GroupId   uint   `json:"group_id"`
	UserEmail string `json:"user_email"`
}
