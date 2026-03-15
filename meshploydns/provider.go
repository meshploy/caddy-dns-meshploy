package meshploydns

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

// Provider implements libdns interfaces for Meshploy's flat-file CoreDNS setup.
type Provider struct {
	ZoneFilePath string
}

func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	// 1. Read the whole file
	data, err := os.ReadFile(p.ZoneFilePath)
	if err != nil {
		return nil, err
	}

	// 2. Bump the serial number
	newContent := bumpSerial(string(data))

	// 3. Append the new TXT records
	var appended []libdns.Record
	for _, rec := range records {
		rr := rec.RR() 
		line := fmt.Sprintf("@ 60 IN %s \"%s\"\n", rr.Type, rr.Data)
		newContent += line
		appended = append(appended, rec)
	}

	// 4. Overwrite the file with the new serial and records
	if err := os.WriteFile(p.ZoneFilePath, []byte(newContent), 0644); err != nil {
		return nil, err
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
	newContent := bumpSerial(strings.Join(newLines, "\n"))
	err = os.WriteFile(p.ZoneFilePath, []byte(newContent), 0644)
	return records, err
}

func bumpSerial(content string) string {
	re := regexp.MustCompile(`(\d+)(\s*;\s*serial)`)
	return re.ReplaceAllString(content, fmt.Sprintf("%d${2}", time.Now().Unix()))
}