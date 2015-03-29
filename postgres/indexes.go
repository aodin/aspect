package postgres

type IndexMethod string

const (
	Gist  IndexMethod = "gist"
	Gin   IndexMethod = "gin"
	Btree IndexMethod = "btree"
	Hash  IndexMethod = "hash"
)
