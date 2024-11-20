package handler

type idPathIn struct {
	ID int64 `path:"id" json:"id" example:"1" doc:"User ID"`
}

type listIn struct {
	Cursor string `query:"cursor"`
	Limit  int    `query:"limit"`
}
