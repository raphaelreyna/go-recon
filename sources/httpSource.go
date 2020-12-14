package sources

import (
	"os"
	"io"
	"net/http"
	"net/url"
	"github.com/raphaelreyna/recon"
)

const HTTPSrc recon.SourceName = "http_source"

type HTTPSource struct {
	Client *http.Client `json:"=" bson:"-" yaml:"-"`
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	valid := err == nil
	valid = valid && u.Hostname() != ""
	return valid
}

func (hs *HTTPSource) AddFileAs(name, destination string, perm os.FileMode) bool {
	rollback := true
	nf, err := os.OpenFile(destination, os.O_CREATE | os.O_WRONLY, perm)
	if err != nil {
		return false
	}
	defer func() {
		nf.Close()
		if rollback {
			os.Remove(nf.Name())
		}
	}()

	client := hs.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Get(name)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()

	_, err = io.Copy(nf, resp.Body)
	if err == nil {
		rollback = false
	}

	return err == nil
}
