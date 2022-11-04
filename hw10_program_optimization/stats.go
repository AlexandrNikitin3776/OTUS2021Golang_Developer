package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
type users []User

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (users, error) {
	bufferedReader := bufio.NewReader(r)
	result := make(users, 0)

	var user User
	for {
		line, err := bufferedReader.ReadBytes('\n')
		if len(line) == 0 && err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}

		if err = json.Unmarshal(line, &user); err != nil {
			return result, err
		}

		result = append(result, user)
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, "."+domain) {
			foundDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[foundDomain]++
		}
	}
	return result, nil
}
