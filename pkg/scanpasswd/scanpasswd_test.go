package scanpasswd

import (
	"log"
	"testing"

	"github.com/projectdiscovery/gologger"
)

func TestScanpasswd(t *testing.T) {
	list, err := ListProfiles()
	if err != nil {
		log.Fatal(err)
	}
	for i, p := range list {
		if p.Error != nil {
			gologger.Info().Msgf("%2d. %-30s  [err: %v]\n", i+1, p.SSID, p.Error)
		} else if p.Password == "" {
			gologger.Info().Msgf("%2d. %-30s  <open>\n", i+1, p.SSID)
		} else {
			gologger.Info().Msgf("%2d. %-30s  %s\n", i+1, p.SSID, p.Password)
		}
	}
}
