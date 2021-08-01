package servicex

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/go-courier/docker"
	"github.com/go-courier/helmx/kubetypes"
	"github.com/go-courier/helmx/spec"
	"gopkg.in/yaml.v2"
)

func (conf *Configuration) dockerize() {
	writeToYamlFile("./config/default.yml", conf.defaultConfig(), yaml.Marshal)
	writeToYamlFile("./dockerfile.default.yml", conf.dockerfile(), yaml.Marshal)
	writeToYamlFile("./helmx.default.yml", conf.helmx(), yaml.Marshal)
}

const (
	FromImage = "env-golang"
)

func (conf *Configuration) defaultConfig() map[string]string {
	m := map[string]string{}
	m["GOENV"] = "DEV"

	for _, envVar := range conf.defaultEnvVars.Values {
		if !envVar.Optional {
			m[envVar.Key(conf.Prefix())] = envVar.Value
		}
	}

	return m
}

func (conf *Configuration) helmx() spec.Spec {
	ss := spec.Spec{}
	ss.Service = &spec.Service{}

	keyPrefix := conf.Prefix()

	for _, envVar := range conf.defaultEnvVars.Values {
		v := envVar.Value

		if v != "" {
			if envVar.IsExpose {
				port, err := spec.ParsePort(v)
				if err == nil {
					ss.Service.Ports = append(ss.Service.Ports, *port)
				}
			}

			if envVar.IsHealthCheck {
				action, err := spec.ParseAction(v)
				if err == nil {
					ss.Service.LivenessProbe = &spec.Probe{
						Action: *action,
						ProbeOpts: kubetypes.ProbeOpts{
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
					}
					ss.Service.ReadinessProbe = &spec.Probe{
						Action: *action,
						ProbeOpts: kubetypes.ProbeOpts{
							InitialDelaySeconds: 5,
							PeriodSeconds:       5,
						},
					}
				}
			}
		}

		if envVar.IsUpstream {
			ss.Upstreams = append(ss.Upstreams, "${"+envVar.Key(keyPrefix)+"}")
		}
	}

	if len(ss.Upstreams) > 0 {
		sort.Strings(ss.Upstreams)
	}

	return ss
}

func (conf *Configuration) dockerfile() *docker.Dockerfile {
	dockerfile := &docker.Dockerfile{
		From: FromImage,
	}

	dockerfile = dockerfile.AddEnv("GOENV", "DEV")

	for _, envVar := range conf.defaultEnvVars.Values {
		if envVar.Value != "" {
			if envVar.IsCopy {
				dockerfile = dockerfile.AddContent(envVar.Value, "./")
			}
			if envVar.IsExpose {
				dockerfile = dockerfile.WithExpose(envVar.Value)
			}
		}
	}

	dockerfile = dockerfile.WithWorkDir("/go/bin")
	dockerfile = dockerfile.WithCmd("./"+conf.ServiceName(), "-c=false")
	dockerfile = dockerfile.AddContent("./"+conf.ServiceName(), "./")

	return dockerfile
}

func writeToYamlFile(filename string, v interface{}, marshal func(v interface{}) ([]byte, error)) error {
	bytes, _ := marshal(v)
	dir := filepath.Dir(filename)
	if dir != "" {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, bytes, os.ModePerm)
}
