package storages_test

import (
	"github.com/codex-team/hawk.cloud-manager/pkg/storage/storages"
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
	expected := storages.YamlConfig{}

	err := yaml.Unmarshal([]byte(YamlStorage_Load_File), &expected)
	if err != nil {
		t.Fatal(err)
	}

	tmpfile, err := ioutil.TempFile("", "yaml_load")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Created temp file %s", tmpfile.Name())

	_, err = tmpfile.WriteString(YamlStorage_Load_File)
	if err != nil {
		t.Fatal(err)
	}


	yamlStorage := storages.NewYamlStorage(tmpfile.Name())
	err = yamlStorage.Load()
	if err != nil {
		t.Fatal(err)
	}

	marsh, err := yaml.Marshal(&yamlStorage.Config)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", string(marsh))

	if diff := deep.Equal(expected, yamlStorage.Config); diff != nil {
		t.Error(diff)
	}
}
