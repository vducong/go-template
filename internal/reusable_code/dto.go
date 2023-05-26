package reusablecode

type ReusableCodeGetByCodeReq struct {
	Code string `json:"code" binding:"required"`
}
