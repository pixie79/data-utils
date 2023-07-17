package data_utils

type CredentialsType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TagsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
