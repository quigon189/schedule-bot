package services

import (
	"cmp"
	"context"
	"core/internal/dto"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
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
	Type       int // 0 - обычная, 1 - нечетная, 2 - четная
	SubjectID  int
	GroupID    int
	TeacherID  int
	AudienceID int
}

type TimeSlot struct {
	Day  int
	Slot int
}

type Schedule struct {
	Days  int
	Slots int
	Grid  [][][]*ScheduleCell // [day][slot][idx]
}

func (s *PlannerService) PlanWeeklySchedule(ctx context.Context, req *dto.PlanScheduleRequest) (*dto.WeeklySchedule, error) {
	lr, err := s.generateLessonRequests(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("generate lesson requests: %w")
	}

	prettyLR, _ := json.MarshalIndent(lr, "", "  ")
	log.Printf("Lesson requests:\n%s", prettyLR)

	resp := dto.WeeklySchedule{
		Seed: req.Seed,
	}

	return &resp, nil
}

func (s *PlannerService) generateLessonRequests(ctx context.Context, req *dto.PlanScheduleRequest) ([]LessonRequest, error) {
	lessonRequests := []LessonRequest{}

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
				AudienceID: audience.ID,
				TeacherID:  teacher.User.ID,
				Priority:   subj.Priority,
				IsSplit:    false,
			})
		}

		if subj.LessonsCount%2 != 0 {
			lessonRequests = append(lessonRequests, LessonRequest{
				SubjectID:  subject.ID,
				AudienceID: audience.ID,
				TeacherID:  teacher.User.ID,
				Priority:   subj.Priority,
				IsSplit:    true,
			})
		}
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

func (s *PlannerService) getTimeSlots(sch *Schedule, req *LessonRequest) ([]TimeSlot, error) {
	var timeSlots []TimeSlot
	for slot := 1; slot <= 5; slot++ {
		for day := 1; day <= 5; day++ {
			schCells := sch.Grid[day][slot]
			if !req.IsSplit {
				if slices.ContainsFunc(schCells, func(cell *ScheduleCell) bool {
					if cell.GroupID == req.GroupID ||
						cell.TeacherID == req.TeacherID ||
						cell.AudienceID == req.AudienceID {
						return true
					}
					return false
				}) {
					continue
				}
			}

			timeSlots = append(timeSlots, TimeSlot{
				Day:  day,
				Slot: slot,
			})
		}
	}
	for slot := 1; slot <= 4; slot++ {
		timeSlots = append(timeSlots, TimeSlot{
			Day:  6,
			Slot: slot,
		})
	}

	return timeSlots, nil
}
