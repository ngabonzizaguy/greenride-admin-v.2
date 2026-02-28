package services

import "testing"

func TestNormalizeSMSPhone(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		ok     bool
		output string
	}{
		{
			name:   "rwanda e164 with dash",
			input:  "+250-784928786",
			ok:     true,
			output: "+250784928786",
		},
		{
			name:   "rwanda local leading zero",
			input:  "0784928786",
			ok:     true,
			output: "+250784928786",
		},
		{
			name:   "rwanda no plus",
			input:  "250784928786",
			ok:     true,
			output: "+250784928786",
		},
		{
			name:  "rwanda accidental trunk zero after country",
			input: "+250-0784928786",
			ok:    false,
		},
		{
			name:  "too short invalid",
			input: "+250-88510542",
			ok:    false,
		},
		{
			name:  "non digit invalid",
			input: "+25078A928786",
			ok:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := normalizeSMSPhone(tt.input)
			if ok != tt.ok {
				t.Fatalf("ok mismatch: got=%v want=%v (got value=%q)", ok, tt.ok, got)
			}
			if tt.ok && got != tt.output {
				t.Fatalf("normalized phone mismatch: got=%q want=%q", got, tt.output)
			}
		})
	}
}

