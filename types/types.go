// Description: A package containing useful reusable snippets of GO code
// Author: Pixie79
// ============================================================================
// package data_utils

package types

// CredentialsType is a struct to hold credentials
type CredentialsType struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TagsType is a struct to hold tags
type TagsType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
