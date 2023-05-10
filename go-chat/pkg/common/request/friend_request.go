package request

type FriendRequest struct {
	Uuid           string `json:"uuid"`
	FriendUsername string `json:"friendUsername"`
}
