package services

import (
	"bytes"
	"context"
	"core/internal/models"
	"fmt"
	"text/template"
)

type ExportFormat string

const (
	FormatHTML ExportFormat = "html"
	FormatPDF  ExportFormat = "pdf"
	FormatPNG  ExportFormat = "png"
)

type ScheduleExportService struct {
	scheduleService *ScheduleService
}

func NewScheduleExportService(scheduleService *ScheduleService) *ScheduleExportService {
	return &ScheduleExportService{scheduleService: scheduleService}
}

type GroupScheduleData struct {
	GroupName      string
	AcademicPeriod string
	Days           []Day
}

type Day struct {
	DayName string
	Lessons []Lesson
}

type Lesson struct {
	LessonNumber string
	IsSplit      bool
	Check        bool
	Details      LessonDetails
	Odd          LessonDetails
	Even         LessonDetails
}

type LessonDetails struct {
	Subject  string
	Teacher  string
	Audience string
}

func (s *ScheduleExportService) ExportGroupSchedule(ctx context.Context, groupID int, periodID *int, format ExportFormat) ([]byte, string, error) {
	var period *models.AcademicPeriod
	var err error
	if periodID != nil {
		period, err = s.scheduleService.GetAcademicPeriodByID(ctx, *periodID)
		if err != nil {
			return nil, "", fmt.Errorf("get academic period: %w", err)
		}
	} else {
		period, err = s.scheduleService.GetActiveAcademicPeriod(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("get active academic period: %w", err)
		}
	}

	group, err := s.scheduleService.GetGroup(ctx, groupID)
	if err != nil {
		return nil, "", fmt.Errorf("get group: %w", err)
	}

	templates, err := s.scheduleService.GetGroupSchedule(ctx, groupID, &period.ID)
	if err != nil {
		return nil, "", fmt.Errorf("get group schedule: %w", err)
	}

	data := buildScheduleDate(group, period, templates)

	htmlContent, err := renderScheduleHTML(data)
	if err != nil {
		return nil, "", fmt.Errorf("render HTML: %w", err)
	}

	return htmlContent, "text/html", nil
}

func buildScheduleDate(group *models.Group, period *models.AcademicPeriod, templates []models.ScheduleTemplate) GroupScheduleData {
	data := GroupScheduleData{
		GroupName:      group.Name,
		AcademicPeriod: fmt.Sprintf("%d семестр %s уч.г.", period.Semester, period.Year),
	}

	scheduleMatrix := make(map[int]map[int]*Lesson)
	for day := 1; day <= 7; day++ {
		scheduleMatrix[day] = make(map[int]*Lesson)
		for number := 0; number <= 8; number++ {
			scheduleMatrix[day][number] = &Lesson{}
		}
	}

	for _, tmpl := range templates {
		details := LessonDetails{
			Subject:  tmpl.Subject.Title,
			Teacher:  tmpl.Teacher.User.GetShortName(),
			Audience: tmpl.Audience.Number,
		}
		lesson := scheduleMatrix[tmpl.DayOfWeek][tmpl.Number]
		lesson.LessonNumber = fmt.Sprintf("%d", tmpl.Number)
		lesson.Check = true
		switch tmpl.WeekType {
		case 0:
			lesson.IsSplit = false
			lesson.Details = details
		case 1:
			lesson.IsSplit = true
			lesson.Odd = details
		case 2:
			lesson.IsSplit = true
			lesson.Even = details
		}
	}

	daysOfWeek := map[int]string{
		1: "Понедельник",
		2: "Вторник",
		3: "Среда",
		4: "Четверг",
		5: "Пятница",
		6: "Суббота",
		7: "Воскресенье",
	}

	for d := 1; d <= 6; d++ {
		day := Day{
			DayName: daysOfWeek[d],
			Lessons: []Lesson{},
		}
		for number := 0; number <= 8; number++ {
			if scheduleMatrix[d][number] != nil {
				lesson := scheduleMatrix[d][number]
				if lesson.Check {
					day.Lessons = append(day.Lessons, *lesson)
				}
			}
		}

		if len(day.Lessons) == 0 {
			Lesson := Lesson{
				Details: LessonDetails{
					Subject: "День самостоятельной работы",
				},
			}
			day.Lessons = append(day.Lessons, Lesson)
		}
		data.Days = append(data.Days, day)
	}

	return data
}

func renderScheduleHTML(data GroupScheduleData) ([]byte, error) {
	tmpl, err := template.New("schedule").Parse(tmplStr)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const tmplStr = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>Расписание {{.GroupName}}</title>
    <style>
        /* Настройки страницы для экспорта */
        @page {
            size: A4;
            margin: 0;
        }

        body {
            margin: 0;
            padding: 0;
            font-family: 'Segoe UI', Roboto, Arial, sans-serif;
            background-color: #ffffff;
        }

        .page {
            width: 210mm;
            min-height: 297mm;
            padding: 15mm;
            margin: 0 auto;
            box-sizing: border-box;
            background: white;
        }

        /* Заголовок */
        .header {
            text-align: center;
            margin-bottom: 20px;
        }

        .header h1 {
            margin: 0;
            font-size: 22pt;
            color: #333;
        }

        .header p {
            margin: 5px 0;
            font-size: 14pt;
            color: #666;
        }

        /* Стили дня недели */
        .day-title {
            font-weight: bold;
            font-size: 16pt;
            margin-top: 20px;
            margin-bottom: 8px;
            color: #2c3e50;
            border-bottom: 2px solid #eee; /* Тонкая подчеркивающая линия для структуры */
        }

        /* Таблица занятий */
        .schedule-table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 10px;
        }

        .schedule-table td {
            padding: 10px;
            border: 1px solid #dee2e6;
            font-size: 11pt;
            vertical-align: middle;
        }

        /* Ширина колонок */
        .col-number { width: 5%; text-align: center; font-weight: bold; background: #f8f9fa; }
        .col-subject { width: 55%; }
        .col-teacher { width: 30%; font-style: italic; }
        .col-room    { width: 10%; text-align: center; }

        /* Чтобы избежать разрыва таблицы при печати в PDF */
        .day-container {
            page-break-inside: avoid;
            margin-bottom: 15px;
        }

		/* Добавляем стили для четности */
    .week-type {
        font-size: 8pt;
        color: #888;
        text-transform: uppercase;
        display: block;
        margin-bottom: 2px;
    }
    
    /* Выделение строк для визуального разделения внутри одной пары */
    .row-odd { background-color: #fff; }
    .row-even { background-color: #fafafa; }
    
    .schedule-table td {
        line-height: 1.2;
    }
    </style>
</head>
<body>

<div class="page">
    <div class="header">
        <h1>Расписание группы {{.GroupName}}</h1>
        <p>{{.AcademicPeriod}}</p>
    </div>

	{{range .Days}}
	<div class="day-container">
		<div class="day-title">{{.DayName}}</div>
		<table class="schedule-table">
			{{range .Lessons}}
				{{if .IsSplit}}
					<!-- Вариант: Разные предметы для Четной/Нечетной -->
					<tr class="row-odd">
						<td class="col-number" rowspan="2">{{.LessonNumber}}</td>
						<td class="col-subject"><span class="week-type">Нечетная</span>{{.Odd.Subject}}</td>
						<td class="col-teacher">{{.Odd.Teacher}}</td>
						<td class="col-room">{{.Odd.Audience}}</td>
					</tr>
					<tr class="row-even">
						<!-- Номер пары пропущен, так как он объединен (rowspan) -->
						<td class="col-subject"><span class="week-type">Четная</span>{{.Even.Subject}}</td>
						<td class="col-teacher">{{.Even.Teacher}}</td>
						<td class="col-room">{{.Even.Audience}}</td>
					</tr>
				{{else}}
					<!-- Вариант: Обычная пара -->
					<tr>
						<td class="col-number">{{.LessonNumber}}</td>
						<td class="col-subject">{{.Details.Subject}}</td>
						<td class="col-teacher">{{.Details.Teacher}}</td>
						<td class="col-room">{{.Details.Audience}}</td>
					</tr>
				{{end}}
			{{end}}
		</table>
	</div>
    {{end}}
</div>

</body>
</html>`
