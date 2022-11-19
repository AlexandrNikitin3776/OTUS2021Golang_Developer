package hw10programoptimization

import (
	"bufio"
	"github.com/mailru/easyjson"
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

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	errChannel := make(chan error)
	defer close(errChannel)

	usersChannel := getUsers(r, errChannel)
	return countDomains(usersChannel, errChannel, domain)
}

func getUsers(r io.Reader, errChannel chan error) <-chan User {
	usersChannel := make(chan User)

	go func() {
		defer close(usersChannel)

		var user User

		bufferedReader := bufio.NewReader(r)
		for {
			line, err := bufferedReader.ReadBytes('\n')
			if len(line) == 0 && err != nil {
				if err != io.EOF {
					errChannel <- err
				}
				break
			}

			if err = easyjson.Unmarshal(line, &user); err != nil {
				errChannel <- err
				break
			}

			usersChannel <- user
		}
	}()

	return usersChannel
}

func countDomains(usersChannel <-chan User, errChannel <-chan error, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for {
		select {
		case user, ok := <-usersChannel:
			if !ok {
				return result, nil
			}
			if strings.Contains(user.Email, "."+domain) {
				foundDomain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
				result[foundDomain]++
			}
		case err := <-errChannel:
			return nil, err
		}
	}
}
