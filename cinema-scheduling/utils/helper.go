package utils

// ---------------- String to Pointer ----------------
// Converts a string to *string, returns nil if empty
func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
