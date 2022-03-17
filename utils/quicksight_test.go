package utils

import (
	"fmt"
	"testing"
)

func TestInitClient(t *testing.T) {
	client := InitClient(nil)
	if client == nil {
		t.Fatal("init client failed")
	}
	region := "us-west-2"
	client = InitClient(&region)
	if client == nil {
		t.Fatal("init client failed")
	}

}

func TestListDashboards(t *testing.T) {
	if accountId, err := GetAccountId(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(*accountId)
		region := "us-west-2"
		client := InitClient(&region)
		dashboards, err := ListDashboards(*accountId, client)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range dashboards {
			t.Log(v)
		}
	}

}

func TestGetAccountId(t *testing.T) {
	if accountId, err := GetAccountId(); err != nil {
		t.Fatal(err)
	} else {
		t.Log(*accountId)
	}

}

func TestGetUsers(t *testing.T) {
	client := InitClient(nil)
	awsAccountId, err := GetAccountId()
	if err != nil {
		t.Fatal(err)
	}

	users := GetUsers(*awsAccountId, "QSDER", client)

	for _, user := range users {
		t.Log(*user.UserName)
	}

}

func TestDeleteUser(t *testing.T) {
	client := InitClient(nil)
	awsAccountId, err := GetAccountId()
	if err != nil {
		t.Fatal(err)
	}
	userName := "QSDER/business@demo.com" //在IAM的用户是由"{role}/{email}"组成
	DeleteQSUser(*awsAccountId, userName, client)
}

func TestCreateUser(t *testing.T) {
	client := InitClient(nil)
	awsAccountId, err := GetAccountId()
	if err != nil {
		t.Fatal(err)
	}

	err = CreateQSUser(*awsAccountId, "QSDER", "business@demo.com", client)
	if err != nil {
		t.Fatal(err)
	}
}
func TestGetEmbedUrl(t *testing.T) {

	awsAccountId, err := GetAccountId()
	if err != nil {
		t.Fatal(err)
	}

	region := "us-west-2"
	roleName := "QSDER"
	userEmail := "business@demo.com"
	dashboardID := "5558d74f-1a50-40d3-8a59-959f87156b4e"

	url := GetEmbedUrl(*awsAccountId, region, roleName, dashboardID, userEmail)
	fmt.Println(*url)

}
