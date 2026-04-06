package llm

import (
	"context"
	"core/internal/dto"
	"core/internal/services"
	"fmt"
	"log"
	"strings"
)

type Tool struct {
	Name        string                                                           `json:"name"`
	Description string                                                           `json:"description"`
	Parameters  map[string]string                                                `json:"parameters"`
	Handler     func(ctx context.Context, params map[string]any) (string, error) `json:"-"`
}

func initTools(svc services.ScheduleService) []Tool {
	return []Tool{
		{
			Name:        "get_schedule_group",
			Description: "Получить основное расписание группы (без изменений) на неделю для указаного учебного периода. Можешь использовать его, если в запросе указаны на понедельник, на вторник и т.д. до конца недели; а так же на неделю и похожие запросы",
			Parameters: map[string]string{
				"group_name":      "string, обязательный, трехзначиное число (например 501)",
				"period_year":     "string, опциональный, учебный период указывающий года обучения в формате YYYY/YYYY (например 2025/2026); если не указан, то подставляется текущий активный период",
				"period_semester": "int, опциональный (должен быть задан, если задан period_year), указывает семестр учебного года, принимает одно из двух значений: 1 или 2",
			},
			Handler: func(ctx context.Context, params map[string]any) (string, error) {
				log.Printf("Параметры: %+v", params)
				groupName, ok := params["group_name"].(string)
				if !ok {
					return "", fmt.Errorf("group_name обязательный параметр")
				}

				var periodID *int
				periodYear, ok := params["period_year"].(string)
				if ok && periodYear != "" {
					semester, okk := params["period_semester"].(float64)
					if !okk {
						return "", fmt.Errorf("period_semester обязательный, если указан period_year")
					}

					period, err := svc.GetAcademicPeriodByYear(ctx, periodYear, int(semester))
					if err != nil {
						return "", fmt.Errorf("get academic period: %w", err)
					}

					periodID = &period.ID
				}

				group, err := svc.GetGroupByName(ctx, groupName)
				if err != nil {
					return "", fmt.Errorf("get group by name: %w", err)
				}

				schedule, err := svc.GetAllScheduleTemplates(ctx, &dto.ScheduleFiltersRequest{
					GroupID:          &group.ID,
					AcademicPeriodID: periodID,
				})

				var builder strings.Builder
				builder.WriteString("### Основное распиание занятий для группы " + group.Name + ":\n\n")

				dayOfWeek := map[int]string{
					1: "Понедельник",
					2: "Вторник",
					3: "Среда",
					4: "Четверг",
					5: "Пятница",
					6: "Суббота",
					7: "Воскресенье",
				}

				for _, s := range schedule {
					fmt.Fprintf(&builder, "[%s] Пара №%d\n", dayOfWeek[s.DayOfWeek], s.Number)
					fmt.Fprintf(&builder, "- Предмет: %s\n", s.Subject.Title)
					fmt.Fprintf(&builder, "- Аудитория: %s\n", s.Audience.Number)
					fmt.Fprintf(&builder, "- Преподаватель: %s\n", s.Teacher.User.FullName)
					builder.WriteString("\n")
				}

				return builder.String(), nil
			},
		},
	}
}
