package sf

type Service struct {
	Name                 string                `yaml:"name"`
	Src                  string                `yaml:"src"`
	Replicas             int32                 `yaml:"replicas"`
	EnvironmentVariables []EnvironmentVariable `yaml:"environment"`
	AffinityLabels       map[string]string     `yaml:"affinityLabels"`
	Volumes              []Volume              `yaml:"volumes"`
	Setup                Setup                 `yaml:"setup"`
	Config               Config                `yaml:"config"`
	Data                 Data                  `yaml:"data"`
}
type Data struct {
	File        string `yaml:"file"`
	Destination string `yaml:"destination"`
}
type Config struct {
	Script      string `yaml:"script"`
	Destination string `yaml:"destination"`
}
type Setup struct {
	Script string `yaml:"script"`
	User   string `yaml:"user"`
}
type Volume struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
	Driver      string `yaml:"driver"`
}
type Port struct {
	Port     int32  `yaml:"port"`
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
}
type EnvironmentVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
type Endpoint struct {
	Name    string `yaml:"name"`
	Service string `yaml:"service"`
	Ports   []Port `yaml:"ports"`
}
type Ingress struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}
type Rule struct {
	Protocol string `yaml:"protocol"`
	Endpoint string `yaml:"endpoint"`
}
type Mesh struct {
	Endpoints []Endpoint `yaml:"endpoints"`
	Ingresses []Ingress  `yaml:"ingresses"`
}
type Egg struct {
	Name     string            `yaml:"name"`
	Services []Service         `yaml:"services"`
	Mesh     Mesh              `yaml:"mesh"`
	Metadata map[string]string `yaml:"metadata"`
}
