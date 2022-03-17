package utils

import "testing"

func TestLoadUserFromFile(t *testing.T) {
	fileName := "users.yaml"
	users, err := LoadUsersFromFile(fileName)

	if err != nil {
		t.Fatal(err)
	}

	for _, v := range users {
		t.Log(v)
	}
}

func TestSha1(t *testing.T) {
	password := Sha1("Passw0rd")
	t.Log(password)
}

func TestWriteUserToFile(t *testing.T) {
	data := []User{}
	data = append(data, User{
		Email:      "business@demo.com",
		Password:   Sha1("Passw0rd"),
		DashBoards: []string{"1111111", "2222222"},
	})
	data = append(data, User{
		Email:      "it@demo.com",
		Password:   Sha1("Passw0rd"),
		DashBoards: []string{"33333333", "444444"},
	})

	WriteUserToFile("./test.yaml", data)
}
