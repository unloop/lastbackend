//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2020] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package config

const (
	DefaultBindServerAddress = "0.0.0.0"
	DefaultBindServerPort    = 2992
	DefaultRootDir           = "/var/lib/lastbackend/"
	DefaultCIDR              = "172.0.0.0/24"
)

type Config struct {
	Debug         bool           `yaml:"debug"`
	RootDir       string         `yaml:"root-dir"`
	StorageDriver string         `yaml:"storage-driver"`
	ManifestDir   string         `yaml:"manifest-dir"`
	CIDR          string         `yaml:"cidr"`
	Security      SecurityConfig `yaml:"security"`
	Server        ServerConfig   `yaml:"server"`
	API           NodeClient     `yaml:"api"`
}

type SecurityConfig struct {
	Token string `yaml:"token"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port uint   `yaml:"port"`
	TLS  struct {
		Verify   bool   `yaml:"verify"`
		FileCA   string `yaml:"ca"`
		FileCert string `yaml:"cert"`
		FileKey  string `yaml:"key"`
	} `yaml:"tls"`
}

type NodeClient struct {
	Address string `yaml:"uri"`
	TLS     struct {
		Verify   bool   `yaml:"verify"`
		FileCA   string `yaml:"ca"`
		FileCert string `yaml:"cert"`
		FileKey  string `yaml:"key"`
	} `yaml:"tls"`
}
