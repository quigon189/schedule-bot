package services

import (
	"cmp"
	"context"
	"core/internal/dto"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
)

type WeekType int

const (
	EveryWeek WeekType = 0
	OddWeek   WeekType = 1
	EvenWeek  WeekType = 2
)

type PlannerService struct {
	scheduleService *ScheduleService
}

func NewPlannerService(svc *ScheduleService) *PlannerService {
	return &PlannerService{scheduleService: svc}
}

var seed uint64 = 123

type LessonRequest struct {
	SubjectID  int
	GroupID    int
	TeacherID  int
	AudienceID int
	Priority   int
	seed       uint64
	IsSplit    bool
}

type ScheduleCell struct {
	Type       WeekType
	SubjectID  int
	GroupID    int
	TeacherID  int
	AudienceID int
}

type TimeSlot struct {
	Day  int
	Slot int
	Type WeekType
}

type Schedule struct {
	Days  int
	Slots int
	Grid  [][][]*ScheduleCell // [day][slot][idx]
}

func (s *PlannerService) PlanWeeklySchedule(ctx context.Context, req *dto.PlanScheduleRequest) (*dto.WeeklySchedule, error) {
	lr, err := s.generateLessonRequests(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("generate lesson requests: %w", err)
	}

	log.Printf("Запросов на пары: %d", len(lr))

	if req.Days == 0 {
		req.Days = 5
	}
	if req.Slots == 0 {
		req.Slots = 5
	}

	schedule := Schedule{
		Days:  req.Days,
		Slots: req.Slots,
	}
	schedule.Grid = make([][][]*ScheduleCell, schedule.Days)

	for d := 0; d < schedule.Days; d++ {
		schedule.Grid[d] = make([][]*ScheduleCell, schedule.Slots)
	}

	if err := planSchedule(&schedule, lr, 0, req.Days, req.Slots); err != nil {
		return nil, fmt.Errorf("plan schedule: %w", err)
	}

	respGrid := make(map[int]map[int][]dto.ScheduleCell)

	for day := 1; day <= schedule.Days; day++ {
		respGrid[day] = make(map[int][]dto.ScheduleCell)
		for slot := 1; slot <= schedule.Slots; slot++ {
			respGrid[day][slot] = make([]dto.ScheduleCell, 0)
			for _, cell := range schedule.Grid[day-1][slot-1] {
				subject, _ := s.scheduleService.GetSubjectByID(ctx, cell.SubjectID)
				teacher, _ := s.scheduleService.GetTeacher(ctx, cell.TeacherID)
				audience, _ := s.scheduleService.GetAudience(ctx, cell.AudienceID)

				respGrid[day][slot] = append(respGrid[day][slot], dto.ScheduleCell{
					WeekType: int(cell.Type),
					Subject:  *subject,
					Teacher:  *teacher,
					Audience: *audience,
				})
			}
		}
	}

	resp := dto.WeeklySchedule{
		Seed: req.Seed,
		Grid: respGrid,
	}

	return &resp, nil
}

func (s *PlannerService) generateLessonRequests(ctx context.Context, req *dto.PlanScheduleRequest) ([]LessonRequest, error) {
	lessonRequests := []LessonRequest{}

	teacherBusy := make(map[int]int)

	for _, subj := range req.SubjectList {
		subject, err := s.scheduleService.GetSubjectByID(ctx, subj.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("get subject with id %d: %w", subj.SubjectID, err)
		}
		if subject == nil {
			return nil, fmt.Errorf("subject %d not found", subj.SubjectID)
		}

		teacher, err := s.scheduleService.GetTeacher(ctx, subj.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("get teacher with id %d: %w", subj.TeacherID, err)
		}
		if teacher == nil {
			return nil, fmt.Errorf("teacher %d not found", subj.TeacherID)
		}

		_, ok := teacherBusy[teacher.User.ID]
		if !ok {
			teacherBusy[teacher.User.ID] = 1
		} else {
			teacherBusy[teacher.User.ID]++
		}

		audience, err := s.scheduleService.GetAudience(ctx, subj.AudienceID)
		if err != nil {
			return nil, fmt.Errorf("get audience with id %d: %w", subj.AudienceID, err)
		}
		if audience == nil {
			return nil, fmt.Errorf("audience %d not found", subj.AudienceID)
		}

		for i := 0; i < subj.LessonsCount/2; i++ {
			lessonRequests = append(lessonRequests, LessonRequest{
				SubjectID:  subject.ID,
				GroupID:    subject.GroupID,
				AudienceID: audience.ID,
				TeacherID:  teacher.User.ID,
				Priority:   subj.Priority,
				IsSplit:    false,
			})
		}

		if subj.LessonsCount%2 != 0 {
			lessonRequests = append(lessonRequests, LessonRequest{
				SubjectID:  subject.ID,
				GroupID:    subj.SubjectID,
				AudienceID: audience.ID,
				TeacherID:  teacher.User.ID,
				Priority:   subj.Priority,
				IsSplit:    true,
			})
		}
	}

	var maxLessons int

	for _, val := range teacherBusy {
		if val > maxLessons {
			maxLessons = val
		}
	}

	for i := range lessonRequests {
		tCount := teacherBusy[lessonRequests[i].TeacherID]
		tPriority := tCount / maxLessons * 100
		lessonRequests[i].Priority = (lessonRequests[i].Priority + tPriority) / 2
	}

	if req.Seed == 0 {
		req.Seed = rand.IntN(1000000)
	}

	pcg := rand.NewPCG(seed, uint64(req.Seed))

	r := rand.New(pcg)

	for i := range lessonRequests {
		lessonRequests[i].seed = r.Uint64()
	}

	slices.SortFunc(lessonRequests, func(a, b LessonRequest) int {
		if a.Priority != b.Priority {
			return cmp.Compare(b.Priority, a.Priority)
		}

		return cmp.Compare(b.seed, a.seed)
	})

	return lessonRequests, nil
}

func getTimeSlots(sch *Schedule, req *LessonRequest, days, slots int) []TimeSlot {
	var allTimeSlots []TimeSlot
	var resTimeSlots []TimeSlot

	if days < 6 {
		for slot := 0; slot < slots-1; slot++ {
			for day := 0; day < days-1; day++ {
				allTimeSlots = append(allTimeSlots, TimeSlot{
					Day:  day,
					Slot: slot,
				})
			}
		}
	}
	if days == 6 {
		for slot := 0; slot < slots-1; slot++ {
			for day := 0; day < days-2; day++ {
				allTimeSlots = append(allTimeSlots, TimeSlot{
					Day:  day,
					Slot: slot,
				})
			}
		}
		for slot := 0; slot < slots-1; slot++ {
			allTimeSlots = append(allTimeSlots, TimeSlot{
				Day:  4,
				Slot: slot,
			})

		}
	}

	for _, slot := range allTimeSlots {
		if req.IsSplit {
			if checkTimeSlotSchedule(OddWeek, sch.Grid[slot.Day][slot.Slot], req) {
				resTimeSlots = append(resTimeSlots, TimeSlot{
					Day:  slot.Day,
					Slot: slot.Slot,
					Type: OddWeek,
				})
			}
			if checkTimeSlotSchedule(EvenWeek, sch.Grid[slot.Day][slot.Slot], req) {
				resTimeSlots = append(resTimeSlots, TimeSlot{
					Day:  slot.Day,
					Slot: slot.Slot,
					Type: EvenWeek,
				})
			}
		} else {
			if checkTimeSlotSchedule(EveryWeek, sch.Grid[slot.Day][slot.Slot], req) {
				resTimeSlots = append(resTimeSlots, TimeSlot{
					Day:  slot.Day,
					Slot: slot.Slot,
					Type: EveryWeek,
				})
			}
		}
	}

	if req.IsSplit {
		var priority1 []TimeSlot
		var priority2 []TimeSlot
		var other []TimeSlot

		for _, slot := range resTimeSlots {
			var teacherFree bool
			var groupFree bool
			for _, cell := range sch.Grid[slot.Day][slot.Slot] {
				if cell.Type == EveryWeek {
					continue
				}
				if cell.Type != slot.Type {
					if cell.TeacherID == req.TeacherID {
						teacherFree = true
					}
					if cell.GroupID == req.GroupID {
						groupFree = true
					}
				}
			}
			if groupFree && teacherFree {
				priority1 = append(priority1, slot)
				continue
			}
			if groupFree {
				priority2 = append(priority2, slot)
				continue
			}
			other = append(other, slot)
		}

		resTimeSlots = []TimeSlot{}

		resTimeSlots = append(resTimeSlots, priority1...)
		resTimeSlots = append(resTimeSlots, priority2...)
		resTimeSlots = append(resTimeSlots, other...)
	}

	return resTimeSlots
}

func planSchedule(sch *Schedule, req []LessonRequest, idx int, days, slots int) error {
	if idx >= len(req) {
		return nil
	}

	r := req[idx]
	schSlots := getTimeSlots(sch, &r, days, slots)

	log.Printf("Plan %d req %+v", idx, req[idx])
	log.Printf("Getted slots: %+v", schSlots)

	for _, slot := range schSlots {
		cell := ScheduleCell{
			Type:       slot.Type,
			GroupID:    r.GroupID,
			TeacherID:  r.TeacherID,
			AudienceID: r.AudienceID,
			SubjectID:  r.SubjectID,
		}
		sch.Grid[slot.Day][slot.Slot] = append(sch.Grid[slot.Day][slot.Slot], &cell)

		log.Printf("Подставили пару: день %d пара %d предмет %d преподаватель %d аудитория %d",
			slot.Day, slot.Slot, cell.SubjectID, cell.TeacherID, cell.AudienceID)

		if err := planSchedule(sch, req, idx+1, days, slots); err == nil {
			return nil
		}
		log.Printf("Отменили пару: день %d пара %d предмет %d преподаватель %d аудитория %d",
			slot.Day, slot.Slot, cell.SubjectID, cell.TeacherID, cell.AudienceID)

		sch.Grid[slot.Day][slot.Slot] = sch.Grid[slot.Day][slot.Slot][:len(sch.Grid[slot.Day][slot.Slot])-1]
	}

	return fmt.Errorf("not a single slot fit")
}

func checkTimeSlotSchedule(t WeekType, sch []*ScheduleCell, req *LessonRequest) bool {
	switch t {
	case EveryWeek:
		return !slices.ContainsFunc(sch, func(sc *ScheduleCell) bool {
			if sc.GroupID == req.GroupID || sc.AudienceID == req.AudienceID || sc.TeacherID == req.TeacherID {
				return true
			}
			return false
		})
	case OddWeek:
		return !slices.ContainsFunc(sch, func(sc *ScheduleCell) bool {
			if sc.Type == EvenWeek {
				return false
			}
			if sc.GroupID == req.GroupID || sc.AudienceID == req.AudienceID || sc.TeacherID == req.TeacherID {
				return true
			}
			return false
		})
	case EvenWeek:
		return !slices.ContainsFunc(sch, func(sc *ScheduleCell) bool {
			if sc.Type == OddWeek {
				return false
			}
			if sc.GroupID == req.GroupID || sc.AudienceID == req.AudienceID || sc.TeacherID == req.TeacherID {
				return true
			}
			return false
		})
	default:
		return false
	}
}
