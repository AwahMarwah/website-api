package merchant

type ReqPath struct {
	Id string `uri:"id" binding:"required"`
}