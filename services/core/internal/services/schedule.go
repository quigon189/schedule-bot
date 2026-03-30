package services

import (
	"context"
	"core/internal/dto"
	"core/internal/models"
	"core/internal/repository"
	"errors"
	"fmt"
	"slices"
	"time"
)

type ScheduleService struct {
	audienceRepo *repository.AudienceRepo
	groupRepo    *repository.GroupRepo
	teacherRepo  *repository.TeacherRepo
	studentRepo  *repository.StudentRepo
	userRepo     *repository.UserRepo
	roleRepo     *repository.RoleRepo
	subjectRepo  *repository.SubjectRepo
}

func NewScheduleService(
	audienceRepo *repository.AudienceRepo,
	groupRepo *repository.GroupRepo,
	teacherRepo *repository.TeacherRepo,
	studentRepo *repository.StudentRepo,
	userRepo *repository.UserRepo,
	roleRepo *repository.RoleRepo,
	subjectRepo *repository.SubjectRepo,
) *ScheduleService {
	return &ScheduleService{
		audienceRepo: audienceRepo,
		groupRepo:    groupRepo,
		teacherRepo:  teacherRepo,
		studentRepo:  studentRepo,
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		subjectRepo:  subjectRepo,
	}
}

// ---------- Audience methods ----------
func (s *ScheduleService) CreateAudience(ctx context.Context, req *dto.CreateAudienceRequest) (*models.Audience, error) {
	audience := models.Audience{
		Name:   req.Name,
		Number: req.Number,
	}
	err := s.audienceRepo.Create(ctx, &audience)
	return &audience, err
}

func (s *ScheduleService) GetAudience(ctx context.Context, id int) (*models.Audience, error) {
	return s.audienceRepo.Get(ctx, id)
}

func (s *ScheduleService) GetAllAudience(ctx context.Context) ([]models.Audience, error) {
	return s.audienceRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateAudience(ctx context.Context, id int, req *dto.UpdateAudienceRequest) (*models.Audience, error) {
	audience, err := s.GetAudience(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get audience: %w", err)
	}

	if req.Name != "" {
		audience.Name = req.Name
	}
	if req.Number != "" {
		audience.Number = req.Number
	}

	err = s.audienceRepo.Update(ctx, audience)
	return audience, err
}

func (s *ScheduleService) DeleteAudience(ctx context.Context, id int) error {
	return s.audienceRepo.Delete(ctx, id)
}

// ---------- Group methods ----------
func (s *ScheduleService) CreateGroup(ctx context.Context, req *dto.CreateGroupRequest) (*models.Group, error) {
	group := models.Group{
		Name:          req.Name,
		Specialty:     req.Specialty,
		AdmissionYear: req.AdmissionYear,
	}
	err := s.groupRepo.Create(ctx, &group)
	return &group, err
}

func (s *ScheduleService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	return s.groupRepo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAllGroups(ctx context.Context) ([]models.Group, error) {
	return s.groupRepo.GetAll(ctx)
}

func (s *ScheduleService) UpdateGroup(ctx context.Context, id int, req *dto.UpdateGroupRequest) (*models.Group, error) {
	group, err := s.GetGroup(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}

	if req.Name != "" {
		group.Name = req.Name
	}

	if req.Specialty != "" {
		group.Specialty = req.Specialty
	}

	if req.AdmissionYear != 0 {
		group.AdmissionYear = req.AdmissionYear
	}

	err = s.groupRepo.Update(ctx, group)
	return group, err
}

func (s *ScheduleService) DeleteGroup(ctx context.Context, id int) error {
	return s.groupRepo.Delete(ctx, id)
}

func (s *ScheduleService) CreateTeacher(ctx context.Context, req *dto.CreateUserRequest) (*models.Teacher, error) {
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.teacherRepo.CreateTeacher(ctx, &user); err != nil {
		return nil, err
	}
	return &models.Teacher{User: user}, nil
}

func (s *ScheduleService) GetTeacher(ctx context.Context, userID int) (*models.Teacher, error) {
	return s.teacherRepo.GetTeacherByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllTeachers(ctx context.Context) ([]models.Teacher, error) {
	return s.teacherRepo.GetAllTeachers(ctx)
}

func (s *ScheduleService) DeleteTeacher(ctx context.Context, userID int) error {
	return s.teacherRepo.DeleteTeacher(ctx, userID)
}

func (s *ScheduleService) CreateStudent(ctx context.Context, req *dto.CreateStudentRequest) (*models.Student, error) {
	group, err := s.groupRepo.GetByID(ctx, req.GroupID)
	if err != nil {
		return nil, fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("group with id %d not found", req.GroupID)
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user := models.User{
		Name:         req.Username,
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}
	if err := s.userRepo.Create(ctx, &user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := s.studentRepo.CreateStudent(ctx, &user, req.GroupID); err != nil {
		return nil, err
	}
	return &models.Student{User: user, Group: group}, nil
}

func (s *ScheduleService) GetStudent(ctx context.Context, userID int) (*models.Student, error) {
	return s.studentRepo.GetStudentByUserID(ctx, userID)
}

func (s *ScheduleService) GetAllStudents(ctx context.Context) ([]models.Student, error) {
	return s.studentRepo.GetAllStudents(ctx)
}

func (s *ScheduleService) UpdateStudentGroup(ctx context.Context, userID int, groupID int) error {
	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group with id %d not found", groupID)
	}
	return s.studentRepo.UpdateStudentGroup(ctx, userID, groupID)
}

func (s *ScheduleService) DeleteStudent(ctx context.Context, userID int) error {
	return s.studentRepo.DeleteStudent(ctx, userID)
}

func (s *ScheduleService) AssignStudent(ctx context.Context, userID int, groupID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	isTeacher, err := s.teacherRepo.IsTeacher(ctx, userID)
	if err != nil {
		return fmt.Errorf("check teacher: %w", err)
	}
	if isTeacher {
		return errors.New("user is already a teacher, cannot assign student")
	}

	group, err := s.groupRepo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("get group: %w", err)
	}
	if group == nil {
		return errors.New("group not found")
	}

	if err := s.studentRepo.CreateStudent(ctx, user, groupID); err != nil {
		return fmt.Errorf("create student profile: %w", err)
	}

	return nil
}

func (s *ScheduleService) AssignTeacher(ctx context.Context, userID int) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	isStudent, err := s.studentRepo.IsStudent(ctx, userID)
	if err != nil {
		return fmt.Errorf("check student: %w", err)
	}
	if isStudent {
		return errors.New("user is already student, cannot assign teacher")
	}

	if err := s.teacherRepo.CreateTeacher(ctx, user); err != nil {
		return fmt.Errorf("create teacher: %w", err)
	}

	return nil
}

func (s *ScheduleService) AssignRole(ctx context.Context, userID int, roleName string) error {
	if roleName != "admin" && roleName != "manager" {
		return errors.New("role must be 'admin' or 'manager'")
	}

	role, err := s.roleRepo.GetRoleByName(ctx, roleName)
	if err != nil {
		return fmt.Errorf("get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	return s.roleRepo.AssignRoleToUser(ctx, userID, roleName)
}

func (s *ScheduleService) RemoveRole(ctx context.Context, userID int, roleName string) error {
	roles, err := s.roleRepo.GetRoles(ctx)
	if err != nil {
		return err
	}

	if !slices.ContainsFunc(roles, func(r models.Role) bool {
		return r.Name == roleName
	}) {
		return errors.New("invalid role name")
	}

	return s.roleRepo.RemoveRoleFromUser(ctx, userID, roleName)
}

func (s *ScheduleService) RemoveStudent(ctx context.Context, userID int) error {
	isStudent, err := s.studentRepo.IsStudent(ctx, userID)
	if err != nil {
		return fmt.Errorf("check student: %w", err)
	}
	if !isStudent {
		return errors.New("user is not a student")
	}

	if err := s.studentRepo.DeleteStudent(ctx, userID); err != nil {
		return fmt.Errorf("delete student profile: %w", err)
	}

	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, "student"); err != nil {
		return fmt.Errorf("remove student role: %w", err)
	}

	return nil
}

func (s *ScheduleService) RemoveTeacher(ctx context.Context, userID int) error {
	isTeacher, err := s.teacherRepo.IsTeacher(ctx, userID)
	if err != nil {
		return fmt.Errorf("check teacher: %w", err)
	}
	if !isTeacher {
		return errors.New("user is not a teacher")
	}

	if err := s.teacherRepo.DeleteTeacher(ctx, userID); err != nil {
		return fmt.Errorf("delete teacher profile: %w", err)
	}

	if err := s.roleRepo.RemoveRoleFromUser(ctx, userID, "teacher"); err != nil {
		return fmt.Errorf("remove student role: %w", err)
	}

	return nil
}

func (s *ScheduleService) GetUserRoles(ctx context.Context, userID int) ([]models.Role, error) {
	return s.roleRepo.GetUserRoles(ctx, userID)
}

// --------Subject Methods
func (s *ScheduleService) CreateSubject(ctx context.Context, req *dto.CreateSubjectRequest) (*models.Subject, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse start_date: %w", err)
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parse end_date: %w", err)
	}
	if startDate.After(endDate) {
		return nil, errors.New("start_date must be before end_date")
	}
	subject := models.Subject{
		Title: req.Title,
		Semester: req.Semester,
		HoursLoad: req.HoursLoad,
		StartDate: startDate,
		EndDate: endDate,
		GroupID: req.GroupID,
	}
	err = s.subjectRepo.CreateSubject(ctx, &subject)
	return &subject, err
}

func (s *ScheduleService) GetSubjectByID(ctx context.Context, id int) (*models.Subject, error) {
	return s.subjectRepo.GetSubjectByID(ctx, id)
}

func (s *ScheduleService) GetAllSubjects(ctx context.Context, req *dto.PaginatedSubjectsRequest) (*dto.PaginatedSubjectsResponse, error) {
	return s.subjectRepo.GetAllSubjects(ctx, req)
}

func (s *ScheduleService) GetSubjectsByGroupID(ctx context.Context, groupID int, req *dto.PaginatedSubjectsRequest) (*dto.PaginatedSubjectsResponse, error) {
	return s.subjectRepo.GetSubjectsByGroupID(ctx, groupID, req)
}

func (s *ScheduleService) UpdateSubject(ctx context.Context, id int, req *dto.UpdateSubjectRequest) (*models.Subject, error) {
	subject, err := s.GetSubjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get subject: %w", err)
	}


	if req.Title != "" {
		subject.Title = req.Title
	}

	if req.Semester != 0 {
		subject.Semester = req.Semester
	}

	if req.HoursLoad != 0 {
		subject.HoursLoad = req.HoursLoad
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("parse start_date: %w", err)
		}
		subject.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("parse end_date: %w", err)
		}
		subject.EndDate = endDate
	}

	if subject.StartDate.After(subject.EndDate) {
		return nil, errors.New("start_date must be before end_date")
	}

	err = s.subjectRepo.UpdateSubject(ctx, subject)
	return subject, err
}

func (s *ScheduleService) DeleteSubject(ctx context.Context, id int) error {
	return s.subjectRepo.DeleteSubject(ctx, id)
}
