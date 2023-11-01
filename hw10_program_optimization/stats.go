package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	countResults := make(DomainStat, bufio.MaxScanTokenSize/2)
	suffix := "." + domain

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")

		if email != "" {
			if strings.HasSuffix(email, suffix) {
				addr := strings.SplitN(email, "@", 2)
				if len(addr) == 2 {
					resultString := strings.ToLower(addr[1])
					countResults[resultString]++
				}
			}
		}
	}
	return countResults, nil
}
