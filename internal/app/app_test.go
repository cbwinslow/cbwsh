package app

import (
	"testing"

	"github.com/cbwinslow/cbwsh/pkg/config"
	"github.com/cbwinslow/cbwsh/pkg/core"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
			errMsg:  "config is nil",
		},
		{
			name: "valid config",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  1000,
				},
				AI: config.AIConfig{
					MonitoringInterval: 30,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "zero history size is valid",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  0,
				},
				AI: config.AIConfig{
					MonitoringInterval: 30,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "negative history size",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  -1,
				},
				AI: config.AIConfig{
					MonitoringInterval: 30,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 30,
				},
			},
			wantErr: true,
			errMsg:  "shell.history_size must be non-negative",
		},
		{
			name: "zero monitoring interval is valid",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  1000,
				},
				AI: config.AIConfig{
					MonitoringInterval: 0,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 30,
				},
			},
			wantErr: false,
		},
		{
			name: "negative monitoring interval",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  1000,
				},
				AI: config.AIConfig{
					MonitoringInterval: -1,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 30,
				},
			},
			wantErr: true,
			errMsg:  "ai.monitoring_interval must be non-negative",
		},
		{
			name: "zero connect timeout is valid",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  1000,
				},
				AI: config.AIConfig{
					MonitoringInterval: 30,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "negative connect timeout",
			cfg: &config.Config{
				Shell: config.ShellConfig{
					DefaultShell: core.ShellTypeBash,
					HistorySize:  1000,
				},
				AI: config.AIConfig{
					MonitoringInterval: 30,
				},
				SSH: config.SSHConfig{
					ConnectTimeout: -1,
				},
			},
			wantErr: true,
			errMsg:  "ssh.connect_timeout must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("validateConfig() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
