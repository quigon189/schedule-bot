package dto

type CreateAcademicPeriodRequest struct {
	Year      string `json:"year" validate:"required"`
	Semester  int    `json:"semester" validate:"required,min=1,max=2"`
	StartDate string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required,datetime=2006-01-02"`
}

type UpdateAcademicPeriodRequest struct {
	Year      string `json:"year" validate:"required_without_all=Semester StartDate EndDate"`
	Semester  int    `json:"semester" validate:"required_without_all=Year StartDate EndDate,min=1,max=2"`
	StartDate string `json:"start_date" validate:"required_without_all=Year Semester EndDate,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"required_without_all=Year Semester StartDate,datetime=2006-01-02"`
}
