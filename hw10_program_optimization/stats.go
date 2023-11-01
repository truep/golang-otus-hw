package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	countResults := make(DomainStat, bufio.MaxScanTokenSize/2)
	suffix := "." + domain

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return countResults, fmt.Errorf("get users error: %w", err)
		}

		if user.Email != "" {
			if strings.HasSuffix(user.Email, suffix) {
				email := strings.SplitN(user.Email, "@", 2)
				if len(email) == 2 {
					resultString := strings.ToLower(email[1])
					countResults[resultString]++
				}
			}
		}
	}
	return countResults, nil
}
