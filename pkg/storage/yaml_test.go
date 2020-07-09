package storage_test

import (
	"github.com/codex-team/hawk.cloud-manager/pkg/config"
	"github.com/codex-team/hawk.cloud-manager/pkg/storage"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
)

var YamlStorage_Load_File = `
hosts:
  - name: test-host
    public_key: TESTPUB

groups:
  - name: test-group
    hosts:
      - test-host`

func TestYamlStorage_Load(t *testing.T) {
	expected := config.PeerConfig{}

	err := yaml.Unmarshal([]byte(YamlStorage_Load_File), &expected)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Expect:\n%+v", expected)

	tmpfile, err := ioutil.TempFile("", "yaml_load")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created temp file %s", tmpfile.Name())

	_, err = tmpfile.WriteString(YamlStorage_Load_File)
	if err != nil {
		t.Fatal(err)
	}


	yamlStorage := storage.NewYamlStorage(tmpfile.Name())
	err = yamlStorage.Load()
	if err != nil {
		t.Fatal(err)
	}

	marsh, err := yaml.Marshal(yamlStorage.Get())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Got:\n%s", string(marsh))

	if diff := deep.Equal(expected, yamlStorage.Get()); diff != nil {
		t.Error(diff)
	}
}
