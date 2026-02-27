# Student Management Program


2. Methods (các phương thức nghiệp vụ)

# Nhóm CRUD cơ bản — gọi repo nhưng có validate trước:

func (s *StudentService) AddStudent(student *models.Student) error
func (s *StudentService) GetAllStudents() ([]*models.Student, error)
func (s *StudentService) GetStudentByID(id string) (*models.Student, error)
func (s *StudentService) UpdateStudent(student *models.Student) error
func (s *StudentService) DeleteStudent(id string) error

# Nhóm điểm số:
func (s *StudentService) AddOrUpdateGrade(studentID, subject string, score float64) error
func (s *StudentService) CalculateGPA(studentID string) (float64, error)
func (s *StudentService) GetAcademicRank(gpa float64) string

# Nhóm tìm kiếm / lọc:
func (s *StudentService) SearchByName(name string) ([]*models.Student, error)
func (s *StudentService) GetStudentsByClass(class string) ([]*models.Student, error)
func (s *StudentService) IsStudentInClass(studentID, class string) (bool, error)

