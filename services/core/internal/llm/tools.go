package llm

type FunctionDescription struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Parameters  map[string]string `json:"parameters"`
}

func GetAllFunctions() []FunctionDescription {
	return []FunctionDescription{
		{
			Name:        "get_group_schedule",
			Description: "Получить расписание группы на указанный учебный период (или текущий, если period_id не указан)",
			Parameters: map[string]string{
				"group_id":  "int, обязательный, ID учебной группы",
				"period_id": "int, опциональный, ID учебного периода (если не указан, то берется текущий активный)",
			},
		},
		{
			Name:        "get_teacher_schedule",
			Description: "Получить расписание преподавателя",
			Parameters: map[string]string{
				"teacher_id": "int, обязательный",
				"period_id":  "int, опциональный",
			},
		},
		{
			Name:        "get_audience_schedule",
			Description: "Получить расписание аудитории",
			Parameters: map[string]string{
				"audience_id": "int, обязательный",
				"period_id":   "int, опциональный",
			},
		},
		{
			Name:        "get_lesson_logs",
			Description: "Получить журнал занятий с фильтрацией",
			Parameters: map[string]string{
				"group_id":   "int, опциональный",
				"subject_id": "int, опциональный",
				"teacher_id": "int, опциональный",
				"date_from":  "string, опциональный, формат YYYY-MM-DD",
				"date_to":    "string, опциональный",
				"status":     "string, опциональный (planned, completed, canceled, rescheduled)",
			},
		},
		{
			Name:        "get_subjects",
			Description: "Получить список предметов с пагинацией",
			Parameters: map[string]string{
				"page":       "int, опциональный",
				"per_page":   "int, опциональный",
				"sort_by":    "string, опциональный",
				"sort_order": "string, опциональный",
			},
		},
		{
			Name:        "get_groups",
			Description: "Получить список всех групп",
			Parameters:  map[string]string{},
		},
		{
			Name:        "get_teachers",
			Description: "Получить список всех преподавателей",
			Parameters:  map[string]string{},
		},
	}
}
