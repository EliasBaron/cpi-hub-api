package entity

type DatabaseFile struct {
	Users  []*UserEntity  `json:"users"`
	Spaces []*SpaceEntity `json:"spaces"`
}
