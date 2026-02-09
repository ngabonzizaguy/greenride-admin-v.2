package utils

// MaskPhone masks the middle portion of a phone number for privacy.
// Returns format like "+25****456" (first 3 + **** + last 3).
func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}
	if len(phone) < 8 {
		return "****"
	}
	return phone[:3] + "****" + phone[len(phone)-3:]
}
