package recon

import (
	"fmt"
	"os"
	"os/user"

	"github.com/mr-pmillz/gorecon/localio"
)

// runDNSRecon runs dnsrecon for all base domains and writes csv output file
func runDNSRecon(domains []string, outputDir string) error {
	// clone dnsrecon repo to /opt if user is root and $HOME/tools if non-root user
	currentUser, _ := user.Current()
	var toolsDir string
	switch currentUser.Username {
	case "root":
		toolsDir = "/opt"
		if err := localio.GitClone("https://github.com/darkoperator/dnsrecon.git", fmt.Sprintf("%s/dnsrecon", toolsDir)); err != nil {
			return err
		}
	default:
		toolsDir = fmt.Sprintf("%s/tools", currentUser.HomeDir)
		if exists, err := localio.Exists(toolsDir); err == nil && !exists {
			if err = os.MkdirAll(toolsDir, 0750); err != nil {
				return err
			}
		}
		if err := localio.GitClone("https://github.com/darkoperator/dnsrecon.git", fmt.Sprintf("%s/dnsrecon", toolsDir)); err != nil {
			return err
		}
	}

	// install virtualenv if not already installed
	if err := localio.InstallPython3VirtualEnv(); err != nil {
		return err
	}

	// setup dnsrecon virtualenv
	virtualEnvsDir := fmt.Sprintf("%s/pyenv", currentUser.HomeDir)
	if exists, err := localio.Exists(virtualEnvsDir); err == nil && !exists {
		if err = os.MkdirAll(virtualEnvsDir, 0750); err != nil {
			return err
		}
	}
	if err := localio.RunCommandPipeOutput(fmt.Sprintf("virtualenv -p python3 %s/dnsrecon", virtualEnvsDir)); err != nil {
		return err
	}

	if err := localio.RunCommandPipeOutput(fmt.Sprintf("source %s/dnsrecon/bin/activate && cd %s/dnsrecon && python3 -m pip install -r requirements.txt", virtualEnvsDir, toolsDir)); err != nil {
		return err
	}

	for _, domain := range domains {
		cmd := fmt.Sprintf("source %s/dnsrecon/bin/activate && cd %s/dnsrecon && python3 dnsrecon.py -d %s -t std,axfr,bing,crt -c %s/%s-dnsrecon-results.csv", virtualEnvsDir, toolsDir, domain, outputDir, domain)
		if err := localio.RunCommandPipeOutput(cmd); err != nil {
			return err
		}
	}

	return nil
}
