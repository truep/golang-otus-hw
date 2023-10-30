package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
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
			if matched := strings.HasSuffix(user.Email, suffix); matched {
				resultString := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
				num := countResults[resultString]
				num++
				countResults[resultString] = num
			}
		}
	}
	return countResults, nil
}
