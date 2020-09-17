package structs

import "database/sql"

type URL struct{
	URI string
	Short string
}

type ReturnURL struct {
	URL string
}

type BaseHandler struct {
	Db *sql.DB
}