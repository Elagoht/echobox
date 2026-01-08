package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		envPort  string
		envRead  string
		envWrite string
		wantPort string
		wantRead int
		wantWrite int
	}{
		{
			name:      "defaults when no env vars",
			wantPort:  DefaultPort,
			wantRead:  DefaultReadTimeout,
			wantWrite: DefaultWriteTimeout,
		},
		{
			name:     "custom values from env vars",
			envPort:  "8080",
			envRead:  "60",
			envWrite: "90",
			wantPort: "8080",
			wantRead:  60,
			wantWrite: 90,
		},
		{
			name:     "invalid env vars use defaults",
			envPort:  "",
			envRead:  "invalid",
			envWrite: "bad",
			wantPort:  DefaultPort,
			wantRead:  DefaultReadTimeout,
			wantWrite: DefaultWriteTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean environment
			os.Unsetenv("PORT")
			os.Unsetenv("READ_TIMEOUT")
			os.Unsetenv("WRITE_TIMEOUT")

			// Set test environment
			if tt.envPort != "" {
				os.Setenv("PORT", tt.envPort)
			}
			if tt.envRead != "" {
				os.Setenv("READ_TIMEOUT", tt.envRead)
			}
			if tt.envWrite != "" {
				os.Setenv("WRITE_TIMEOUT", tt.envWrite)
			}

			got := Load()

			if got.Port != tt.wantPort {
				t.Errorf("Load().Port = %v, want %v", got.Port, tt.wantPort)
			}
			if got.ReadTimeout != tt.wantRead {
				t.Errorf("Load().ReadTimeout = %v, want %v", got.ReadTimeout, tt.wantRead)
			}
			if got.WriteTimeout != tt.wantWrite {
				t.Errorf("Load().WriteTimeout = %v, want %v", got.WriteTimeout, tt.wantWrite)
			}
		})
	}
}
