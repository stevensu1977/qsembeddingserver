package main

import (
	"flag"
	"fmt"
	"stevensu1977/quicksight/utils"
	"strings"
)

//LoadUser load user from yaml
func LoadUser(accoudId, role, region, datafile string) []utils.User {
	users, err := utils.LoadUsersFromFile(datafile)
	if err != nil {
		panic(err)
	}
	return users
}

//AddUserToAWS create quicksight user
func AddUserToAWS(accoudId, role string, user utils.User) error {
	client := utils.InitClient(nil)

	err := utils.CreateQSUser(accoudId, role, user.Email, client)
	return err
}

//createUser create user to AWS and yaml file include password
func createUser(accountId, role, region, datafile, email, password, dashboardId string) {
	fmt.Println(accountId, role, region, datafile, email, password, strings.Split(dashboardId, ","))

	users := LoadUser(accountId, role, region, datafile)
	user := utils.User{
		Email:      email,
		Password:   utils.Sha1(password),
		DashBoards: strings.Split(dashboardId, ","),
	}

	err := AddUserToAWS(accountId, role, user)
	if err != nil {
		panic(err)
	}
	users = append(users, user)

	utils.WriteUserToFile(datafile, users)

}

//listUserFromYaml list user from yaml file
func listUserFromYaml(datafile string) {
	users, err := utils.LoadUsersFromFile(datafile)
	if err != nil {
		panic(err)
	}
	for _, v := range users {
		fmt.Printf("%+v\n", v)
	}
}

//listUserFormAWS list user form AWS
func listUserFormAWS(accountId, role string) {
	client := utils.InitClient(nil)
	users := utils.GetUsers(accountId, role, client)
	for _, v := range users {
		fmt.Println(*v.UserName)
	}

}

//listDashboards AWS quicksight
func listDashboards(accountId, region string) {
	client := utils.InitClient(&region)
	dashboards, err := utils.ListDashboards(accountId, client)
	if err != nil {
		panic(err)
	}

	for _, v := range dashboards {
		fmt.Println(*v.Name, *v.DashboardId)
	}

}

func main() {

	role := flag.String("role", "", "quicksight role")
	accountId := flag.String("accountId", "", "AccountID")
	region := flag.String("qsRegion", "us-west-1", "QuickSight region")
	datafile := flag.String("datafile", "users.yaml", "data file")
	action := flag.String("action", "", "action: createUser|deleteUser|listUser|listDashboard")

	email := flag.String("email", "", "user email")
	password := flag.String("password", "", "password")
	dashboardId := flag.String("dashboardID", "", "dashboard id")

	flag.Parse()

	if *accountId == "" {
		if defaultAccoutId, err := utils.GetAccountId(); err != nil {
			panic(err)
		} else {
			accountId = defaultAccoutId
		}

	}

	fmt.Println(*accountId)

	switch *action {
	case "createUser":
		fmt.Println("Action: createUser")
		if *role == "" {
			fmt.Println("CreateUser, role must be required!")
			return
		}
		createUser(*accountId, *role, *region, *datafile, *email, *password, *dashboardId)
	case "listUser":
		fmt.Println("Action: listUser")
		if *role == "" {
			listUserFromYaml(*datafile)
			return
		}
		listUserFormAWS(*accountId, *role)
	case "listDashboards":
		fmt.Print("Action: listDashboards")
		listDashboards(*accountId, *region)

	default:
		fmt.Println(*accountId)
		fmt.Println("action: createUser|deleteUser|listUser")
	}
}
