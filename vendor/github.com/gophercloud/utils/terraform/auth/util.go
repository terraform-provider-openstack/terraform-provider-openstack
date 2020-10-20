package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// This is copied directly from Terraform in order to remove a single legacy
// vendor dependency.

func pathOrContents(poc string) (string, bool, error) {
	if len(poc) == 0 {
		return poc, false, nil
	}

	path := poc
	if path[0] == '~' {
		var err error
		path, err = homedir.Expand(path)
		if err != nil {
			return path, true, err
		}
	}

	if _, err := os.Stat(path); err == nil {
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return string(contents), true, err
		}
		return string(contents), true, nil
	}

	return poc, false, nil
}

const uaEnvVar = "TF_APPEND_USER_AGENT"

func terraformUserAgent(version, sdkVersion string) string {
	ua := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io)", version)
	if sdkVersion != "" {
		ua += " " + fmt.Sprintf("Terraform Plugin SDK/%s", sdkVersion)
	}

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}
