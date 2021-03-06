package yaml

import (
	"io/ioutil"
	"testing"

	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

var loadFile = `
hosts:
  - name: test-host
    public_key: TESTPUB
    allowed_ips: ["10.11.0.1/24"]

groups:
  - name: test-group
    hosts:
      - test-host`

// TestLoad tests loading config to YamlStorage
func TestLoad(t *testing.T) {
	expected := config.PeerConfig{}

	err := yaml.Unmarshal([]byte(loadFile), &expected)
	require.Nil(t, err)

	tmpfile, err := ioutil.TempFile("", "yaml_load")
	require.Nil(t, err)

	_, err = tmpfile.WriteString(loadFile)
	require.Nil(t, err)

	yamlStorage := NewYamlStorage(tmpfile.Name())
	err = yamlStorage.Load()
	require.Nil(t, err)

	_, err = yaml.Marshal(yamlStorage.Get())
	require.Nil(t, err)

	if diff := deep.Equal(expected, *yamlStorage.Get()); diff != nil {
		t.Error(diff)
	}
}
