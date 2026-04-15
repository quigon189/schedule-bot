package services

import (
	"bytes"
	"context"
	"core/internal/models"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

type ExportFormat string

const (
	FormatHTML ExportFormat = "html"
	FormatPDF  ExportFormat = "pdf"
	FormatPNG  ExportFormat = "png"
	FormatDOCX ExportFormat = "docx"
)

var daysOfWeek = map[int]string{
	1: "Понедельник",
	2: "Вторник",
	3: "Среда",
	4: "Четверг",
	5: "Пятница",
	6: "Суббота",
	7: "Воскресенье",
}

type ScheduleExportService struct {
	scheduleService *ScheduleService
}

func NewScheduleExportService(scheduleService *ScheduleService) *ScheduleExportService {
	return &ScheduleExportService{scheduleService: scheduleService}
}

type ScheduleData struct {
	Title          string
	AcademicPeriod string
	Days           []Day
	WidthCol0      string
	WidthCol1      string
	WidthCol2      string
	WidthCol3      string
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
	Col1 string
	Col2 string
	Col3 string
}

func (s *ScheduleExportService) ExportTeacherSchedule(ctx context.Context, teacherID int, periodID *int, format ExportFormat) ([]byte, string, error) {
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

	teacher, err := s.scheduleService.GetTeacher(ctx, teacherID)
	if err != nil {
		return nil, "", fmt.Errorf("get teacher: %w", err)
	}
	if teacher == nil {
		return nil, "", fmt.Errorf("teacher %d not found", teacherID)
	}

	templates, err := s.scheduleService.GetTeacherSchedule(ctx, teacherID, periodID)
	if err != nil {
		return nil, "", fmt.Errorf("get teacher schedule: %w", err)
	}

	data := buildTeacherScheduleData(teacher, period, templates)
	return convert(data, format)
}

func (s *ScheduleExportService) ExportAudienceSchedule(ctx context.Context, audienceID int, periodID *int, format ExportFormat) ([]byte, string, error) {
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

	audience, err := s.scheduleService.GetAudience(ctx, audienceID)
	if err != nil {
		return nil, "", fmt.Errorf("get audience: %w", err)
	}
	if audience == nil {
		return nil, "", fmt.Errorf("audience %d not found", audienceID)
	}

	templates, err := s.scheduleService.GetAudienceSchedule(ctx, audienceID, periodID)
	if err != nil {
		return nil, "", fmt.Errorf("get audience schedule: %w", err)
	}

	data := buildAudienceScheduleData(audience, period, templates)
	return convert(data, format)
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
	if group == nil {
		return nil, "", fmt.Errorf("group %d not found", groupID)
	}

	templates, err := s.scheduleService.GetGroupSchedule(ctx, groupID, &period.ID)
	if err != nil {
		return nil, "", fmt.Errorf("get group schedule: %w", err)
	}

	data := buildGroupScheduleData(group, period, templates)
	return convert(data, format)
}

func convert(data ScheduleData, format ExportFormat) ([]byte, string, error) {
	htmlContent, err := renderScheduleHTML(data)
	if err != nil {
		return nil, "", fmt.Errorf("render HTML: %w", err)
	}

	if format == FormatHTML {
		return htmlContent, "text/html", nil
	}

	tmpFile, err := os.CreateTemp("", "schedule_*.html")
	if err != nil {
		return nil, "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(htmlContent); err != nil {
		return nil, "", fmt.Errorf("write temp file: %w", err)
	}
	defer tmpFile.Close()

	var out []byte
	var mime string
	switch format {
	case FormatPDF:
		out, err = convertHTMLToPDF(tmpFile.Name())
		mime = "application/pdf"
	case FormatPNG:
		out, err = convertHTMLToPNG(tmpFile.Name())
		mime = "image/png"
	case FormatDOCX:
		out, err = convertHTMLToDOCX(tmpFile.Name())
		mime = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return nil, "", fmt.Errorf("convert to %s: %w", format, err)
	}

	return out, mime, nil
}

func buildTeacherScheduleData(teacher *models.Teacher, period *models.AcademicPeriod, templates []models.ScheduleTemplate) ScheduleData {
	data := ScheduleData{
		Title:          fmt.Sprintf("Расписание преподавателя %s", teacher.User.GetShortName()),
		AcademicPeriod: fmt.Sprintf("%d семестр %s уч.г.", period.Semester, period.Year),
		WidthCol0:      "5%",
		WidthCol1:      "55%",
		WidthCol2:      "20%",
		WidthCol3:      "20%",
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
			Col1: tmpl.Subject.Title,
			Col2: tmpl.Subject.Group.Name,
			Col3: tmpl.Audience.Number,
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
			continue
		}
		data.Days = append(data.Days, day)
	}

	return data
}

func buildAudienceScheduleData(audience *models.Audience, period *models.AcademicPeriod, templates []models.ScheduleTemplate) ScheduleData {
	data := ScheduleData{
		Title:          fmt.Sprintf("Расписание аудитории %s", audience.Name),
		AcademicPeriod: fmt.Sprintf("%d семестр %s уч.г.", period.Semester, period.Year),
		WidthCol0:      "5%",
		WidthCol1:      "55%",
		WidthCol2:      "20%",
		WidthCol3:      "20%",
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
			Col1: tmpl.Subject.Title,
			Col2: tmpl.Subject.Group.Name,
			Col3: tmpl.Teacher.User.GetShortName(),
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
			continue
		}
		data.Days = append(data.Days, day)
	}

	return data
}

func buildGroupScheduleData(group *models.Group, period *models.AcademicPeriod, templates []models.ScheduleTemplate) ScheduleData {
	data := ScheduleData{
		Title:          fmt.Sprintf("Расписание группы %s", group.Name),
		AcademicPeriod: fmt.Sprintf("%d семестр %s уч.г.", period.Semester, period.Year),
		WidthCol0:      "5%",
		WidthCol1:      "55%",
		WidthCol2:      "30%",
		WidthCol3:      "10%",
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
			Col1: tmpl.Subject.Title,
			Col2: tmpl.Teacher.User.GetShortName(),
			Col3: tmpl.Audience.Number,
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
					Col1: "День самостоятельной работы",
				},
			}
			day.Lessons = append(day.Lessons, Lesson)
		}
		data.Days = append(data.Days, day)
	}

	return data
}

func renderScheduleHTML(data ScheduleData) ([]byte, error) {
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

func convertHTMLToPDF(htmlPath string) ([]byte, error) {
	cmd := exec.Command("wkhtmltopdf", "--enable-local-file-access", htmlPath, "-")
	return cmd.Output()
}

func convertHTMLToPNG(htmlPath string) ([]byte, error) {
	cmd := exec.Command("wkhtmltoimage", "--enable-local-file-access", htmlPath, "-")
	return cmd.Output()
}

func convertHTMLToDOCX(htmlPath string) ([]byte, error) {
	cmd := exec.Command("pandoc", htmlPath, "-t", "docx", "-o", "-")
	return cmd.Output()
}

const tmplStr = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
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
        .col-1 { text-align: center; font-weight: bold; background: #f8f9fa; }
		.col-2 { text-align: left; }
		.col-3 { text-align: center; }
        .col-4 { text-align: center; }

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
        <h1>{{.Title}}</h1>
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
						<td width="{{$.WidthCol0}}" class="col-1" rowspan="2">{{.LessonNumber}}</td>
						<td width="{{$.WidthCol1}}" class="col-2"><span class="week-type">Нечетная</span>{{.Odd.Col1}}</td>
						<td width="{{$.WidthCol2}}" class="col-3">{{.Odd.Col2}}</td>
						<td width="{{$.WidthCol3}}" class="col-4">{{.Odd.Col3}}</td>
					</tr>
					<tr class="row-even">
						<!-- Номер пары пропущен, так как он объединен (rowspan) -->
						<td width="{{$.WidthCol1}}" class="col-2"><span class="week-type">Четная</span>{{.Even.Col1}}</td>
						<td width="{{$.WidthCol2}}" class="col-3">{{.Even.Col2}}</td>
						<td width="{{$.WidthCol3}}" class="col-4">{{.Even.Col3}}</td>
					</tr>
				{{else}}
					<!-- Вариант: Обычная пара -->
					<tr>
						<td width="{{$.WidthCol0}}" class="col-1">{{.LessonNumber}}</td>
						<td width="{{$.WidthCol1}}" class="col-2">{{.Details.Col1}}</td>
						<td width="{{$.WidthCol2}}" class="col-3">{{.Details.Col2}}</td>
						<td width="{{$.WidthCol3}}" class="col-4">{{.Details.Col3}}</td>
					</tr>
				{{end}}
			{{end}}
		</table>
	</div>
    {{end}}
</div>

</body>
</html>`
