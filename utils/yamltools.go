package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type User struct {
	Email      string
	Password   string
	DashBoards []string
}

func LoadUsersFromFile(fileName string) ([]User, error) {

	yfile, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	data := make(map[string][]User)

	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {

		return nil, err
	}
	return data["users"], nil

}

func WriteUserToFile(fileName string, users []User) {
	data := make(map[string][]User)
	data["users"] = []User{}

	data["users"] = users

	rawData, err := yaml.Marshal(&data)

	if err != nil {
		panic(err)
	}
	err2 := ioutil.WriteFile(fileName, rawData, 0644)

	if err2 != nil {

		log.Fatal(err2)
	}

}

func Sha1(input string) string {
	h := sha1.New()
	io.WriteString(h, input)
	return hex.EncodeToString(h.Sum(nil))
}
