package models

type Department struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	EnvId int    `json:"env_id"`
}
