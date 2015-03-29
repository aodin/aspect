package postgres

type IndexMethod string

const (
	Gist  IndexMethod = "gist"
	Gin               = "gin"
	Btree             = "btree"
	Hash              = "hash"
)
