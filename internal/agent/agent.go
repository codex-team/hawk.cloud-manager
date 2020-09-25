package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"text/template"

	"github.com/codex-team/hawk.cloud-manager/pkg/api"
)

var toWGConf = `[Interface]
PrivateKey = {{ .PrivKey }}
ListenPort = {{ .Config.ListenPort }}
{{ range .Config.Peers }}
[Peer]
PublicKey = {{ .PublicKey }}
Endpoint = {{ .Endpoint }}
AllowedIPs = {{ formatArray .AllowedIPs }}
{{ end }}
`

// Agent is service running on target host, responsible for applying WireGuard configuration on host
type Agent struct {
	// CloudManager address
	managerAddress string
	// WireGuard config
	Config *api.Conf
	// WireGuard privste key
	PrivKey string
	// WireGuard public key
	PubKey string
	// HTTP client
	client *http.Client
	// File to store config in
	configFile string
}

// New returns new Agent
func New(managerAddress string, configFile string, pubKey string, privKey string) (*Agent, error) {
	return &Agent{
		managerAddress: managerAddress,
		PrivKey:        privKey,
		PubKey:         pubKey,
		client:         &http.Client{},
		configFile:     configFile,
		Config:         new(api.Conf),
	}, nil
}

// parseWGConf parses WireGuard config
func (a *Agent) parseWGConf() (string, error) {
	t, err := template.New("toWGConf").Funcs(template.FuncMap{
		"formatArray": func(ips []string) string {
			return strings.Join(ips, `, `)
		},
	}).Parse(toWGConf)
	if err != nil {
		return "", err
	}
	conf := new(bytes.Buffer)

	err = t.ExecuteTemplate(conf, "toWGConf", *a)
	if err != nil {
		return "", fmt.Errorf("failed to parse config: %w", err)
	}
	return strings.TrimSuffix(conf.String(), "\n\n"), nil
}

// queryConf retrieves latest WireGuard configuration from Manager
func (a *Agent) queryConf() (*api.Conf, error) {
	creds := api.Creds{
		PublicKey: a.PubKey,
	}
	reqBody, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, a.managerAddress+"/topology", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get config: status code %d, error: %s", resp.StatusCode, string(body))
	}

	apiConf := api.Conf{}
	err = json.Unmarshal(body, &apiConf)
	if err != nil {
		return nil, err
	}

	return &apiConf, nil
}

// PullNewConf updates Wireguard configuration on host if it has changed
func (a *Agent) PullNewConf() error {
	conf, err := a.queryConf()
	if err != nil {
		return err
	}

	if !(a.Config.Equals(conf)) {
		a.Config = conf

		// parse config
		wgConf, err := a.parseWGConf()
		if err != nil {
			return err
		}

		// save config to file
		err = a.save(wgConf)
		if err != nil {
			return err
		}

		// apply config
		err = a.apply()
		if err != nil {
			return err
		}
	}

	return nil
}

// save saves WireGuard configuration to file
func (a *Agent) save(conf string) error {
	err := ioutil.WriteFile(a.configFile, []byte(conf), 0644)
	if err != nil {
		return err
	}
	log.Println("written")

	return nil
}

// apply applies WireGuard configuration
func (a *Agent) apply() error {
	cmd := exec.Command("wg", "setconf", "wg0", a.configFile)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply WireGuard configuration: %s", stderr.String())
	}
	return nil
}
