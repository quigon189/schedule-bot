package tools

import "github.com/mark3labs/mcp-go/mcp"

type ScheduleTools struct{}

func NewScheduleTools() *ScheduleTools {
	return &ScheduleTools{}
}

func (t *ScheduleTools) GetChanges() mcp.Tool {
	return mcp.NewTool(
		"get_schedule_changes",
		mcp.WithDescription("Позволяет получить ссылку на изображение(я), в которой указаны изменения в расписании для всех груп на заданную в формате ISO 8601 дату"),
		mcp.WithString("date",
			mcp.Description("Дата в формате ISO 8601 YYYY-MM-DD. Если не задана, то возвращает изменения на текущую дату"),
		),
	)
}

func (t *ScheduleTools) GetCurrentDate() mcp.Tool {
	return mcp.NewTool(
		"get_current_date",
		mcp.WithDescription("Возвращает текущую дату в формате ISO 8601"),
		mcp.WithString("none", mcp.Description("")),
	)
}

func (t *ScheduleTools) GetCurrentStudyPeriod() mcp.Tool {
	return mcp.NewTool(
		"get_current_study_period",
		mcp.WithDescription("Возварщает текущий учебный переод. Текущий учебный год в формате YYYY/YYYY, а также текущее полугодие в формате целого числа 1 или 2"),
		mcp.WithString("none", mcp.Description("")),
	)
}

func (t *ScheduleTools) GetGroupSchedule() mcp.Tool {
	return mcp.NewTool(
		"get_group_schedule",
		mcp.WithDescription("Возвращает расписание указанной группы, на указанный учебный период"),
		mcp.WithString("group_name", mcp.Description("Наименование учебной группы в формате АБ-123")),
		mcp.WithString("academic_year", mcp.Description("Учебнгый год в формате YYYY/YYYY")),
		mcp.WithString("half_year", mcp.Description("Полугодие: одно из двух значений 1 или 2")),
	)
}

func (t *ScheduleTools) SendChanges() mcp.Tool {
	return mcp.NewTool(
		"send_changes_to_user",
		mcp.WithDescription("Отправляет пользователю изменения в расписании"),
		mcp.WithString("date", mcp.Description("Дата в формате ISO 8601 YYYY-MM-DD. Если не задана, то возвращает изменения на текущую дату")),
		mcp.WithString("chat_id", mcp.Description("chat_id пользователя отправившего запрос")),
	)
}
