package meshploydns

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

// Provider implements libdns interfaces for Meshploy's flat-file CoreDNS setup.
type Provider struct {
	ZoneFilePath string
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	f, err := os.OpenFile(p.ZoneFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var appended []libdns.Record
	for _, rec := range records {
        // Extract the raw Resource Record struct
		rr := rec.RR() 

		line := fmt.Sprintf("%s 60 IN %s \"%s\"\n", rr.Name, rr.Type, rr.Data)
		if _, err := f.WriteString(line); err != nil {
			return appended, err
		}
		appended = append(appended, rec)
	}

	time.Sleep(3 * time.Second)
	return appended, nil
}

func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	content, err := os.ReadFile(p.ZoneFilePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	for _, line := range lines {
		keep := true
		for _, rec := range records {
			if strings.Contains(line, rec.RR().Data) {
				keep = false
				break
			}
		}
		if keep {
			newLines = append(newLines, line)
		}
	}

	err = os.WriteFile(p.ZoneFilePath, []byte(strings.Join(newLines, "\n")), 0644)
	return records, err
}