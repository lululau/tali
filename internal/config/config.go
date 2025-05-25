package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// ConfigProfile represents a single profile in the Aliyun CLI config
type ConfigProfile struct {
	Name            string `json:"name"`
	Mode            string `json:"mode"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	RegionID        string `json:"region_id"`
	OssEndpoint     string `json:"oss_endpoint,omitempty"` // Custom field for OSS endpoint
	// Other fields like output_format, language can be added if needed
}

// AliyunConfig represents the structure of ~/.aliyun/config.json
type AliyunConfig struct {
	Current  string          `json:"current"`
	Profiles []ConfigProfile `json:"profiles"`
}

// Config holds the application configuration
type Config struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	OssEndpoint     string
}

// LoadAliyunConfig loads configuration from ~/.aliyun/config.json
func LoadAliyunConfig() (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	configPath := filepath.Join(usr.HomeDir, ".aliyun", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read aliyun config file at %s: %w", configPath, err)
	}

	var config AliyunConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse aliyun config file %s: %w", configPath, err)
	}

	if len(config.Profiles) == 0 {
		return nil, fmt.Errorf("no profiles found in aliyun config file: %s", configPath)
	}

	activeProfileName := config.Current
	if activeProfileName == "" {
		if len(config.Profiles) == 1 {
			activeProfileName = config.Profiles[0].Name
		} else {
			for _, p := range config.Profiles {
				if p.Name == "default" {
					activeProfileName = "default"
					break
				}
			}
			if activeProfileName == "" { // If no "default" profile and current is not set
				return nil, fmt.Errorf("no current profile specified in %s, and no 'default' profile found. Please specify a current profile or name one 'default'", configPath)
			}
		}
	}

	var activeProfile *ConfigProfile
	for i := range config.Profiles {
		if config.Profiles[i].Name == activeProfileName {
			activeProfile = &config.Profiles[i]
			break
		}
	}

	if activeProfile == nil {
		return nil, fmt.Errorf("current profile '%s' not found in aliyun config file: %s", activeProfileName, configPath)
	}

	if activeProfile.AccessKeyID == "" || activeProfile.AccessKeySecret == "" || activeProfile.RegionID == "" {
		return nil, fmt.Errorf("profile '%s' in %s is missing access_key_id, access_key_secret, or region_id", activeProfile.Name, configPath)
	}

	// Resolve OSS Endpoint
	ossEndpoint := activeProfile.OssEndpoint
	if ossEndpoint == "" && activeProfile.RegionID != "" {
		ossEndpoint = fmt.Sprintf("oss-%s.aliyuncs.com", activeProfile.RegionID)
	}

	if ossEndpoint == "" {
		return nil, fmt.Errorf("OSS endpoint could not be determined for profile '%s'. "+
			"Please either set 'oss_endpoint' in your profile in %s, "+
			"or ensure 'region_id' is set for default construction (e.g., oss-<region_id>.aliyuncs.com)",
			activeProfile.Name, configPath)
	}

	return &Config{
		AccessKeyID:     activeProfile.AccessKeyID,
		AccessKeySecret: activeProfile.AccessKeySecret,
		RegionID:        activeProfile.RegionID,
		OssEndpoint:     ossEndpoint,
	}, nil
}
