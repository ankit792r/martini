package flutter

type Session struct {
	ProjectPath string
	Answers     map[CommandID][]string
}

func NewSession(projectPath string) *Session {
	return &Session{
		ProjectPath: projectPath,
		Answers:     make(map[CommandID][]string),
	}
}

func (s *Session) SetAnswers(id CommandID, answers []string) {
	s.Answers[id] = answers
}

func (s *Session) AnswersFor(id CommandID) []string {
	return s.Answers[id]
}
