package utils

func NullString(s *string) string {
	if s == nil || *s == "" {
		return ""
	}
	return *s
}
