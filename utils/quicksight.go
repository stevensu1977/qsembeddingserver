package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/quicksight"
	"github.com/aws/aws-sdk-go/service/sts"
)

const QuickSightNamespace = "default"

//InitClient init AWS session client
func InitClient(useRegion *string) *quicksight.QuickSight {

	// Step 1 - Create a session, AssumeRole
	if useRegion == nil {
		useRegion = aws.String("us-east-1")
	}
	sess, err := session.NewSession(&aws.Config{
		Region: useRegion},
	)
	log.Printf("init client successful , region is %s", *sess.Config.Region)

	if err != nil {
		panic(err)
	}

	client := quicksight.New(sess, &aws.Config{
		Region: useRegion})

	return client
}

//GetAccountId get default account id through STS
func GetAccountId() (*string, error) {

	sess := session.Must(session.NewSession())

	stsClient := sts.New(sess)

	callerIndentityOutput, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}
	log.Printf("AccountId: %s", *callerIndentityOutput.Account)
	return callerIndentityOutput.Account, nil

}

//CreateQSUser create QuickSight User
func CreateQSUser(awsAccountID, roleName, userEmail string, client *quicksight.QuickSight) error {

	// Step 2 - Create QuickSight User
	iamRoleARN := "arn:aws:iam::" + awsAccountID + ":role/" + roleName
	identityType := "IAM"

	userRole := "READER" // default quicksight role is READER

	ruInput := quicksight.RegisterUserInput{
		AwsAccountId: &awsAccountID,
		Email:        &userEmail,
		IamArn:       &iamRoleARN,
		Namespace:    aws.String(QuickSightNamespace),
		IdentityType: &identityType,
		SessionName:  &userEmail,
		UserRole:     &userRole,
	}

	ruOutput, ruOutputError := client.RegisterUser(&ruInput)
	if ruOutputError != nil {
		return ruOutputError

	} else {
		log.Println(ruOutput.String())
	}

	return nil

}

//DeleteQSUser delete QuickSight User
func DeleteQSUser(awsAccountID, userName string, client *quicksight.QuickSight) {

	deleteUserInput := quicksight.DeleteUserInput{
		AwsAccountId: &awsAccountID,
		Namespace:    aws.String(QuickSightNamespace),
		UserName:     &userName,
	}

	deleteUserOutput, err := client.DeleteUser(&deleteUserInput)
	if err != nil {
		panic(err)
	}
	fmt.Println(deleteUserOutput)

}

//GetUsers get users by prefixRole
func GetUsers(awsAccountID, prefixRole string, client *quicksight.QuickSight) []*quicksight.User {

	listUserInput := quicksight.ListUsersInput{
		AwsAccountId: &awsAccountID,
		Namespace:    aws.String(QuickSightNamespace),
	}

	listUserOuput, err := client.ListUsers(&listUserInput)
	if err != nil {
		panic(err)
	}

	users := []*quicksight.User{}

	for index, user := range listUserOuput.UserList {
		if prefixRole != "*" {
			if strings.HasPrefix(*user.UserName, prefixRole) {
				fmt.Println(index, user)
				users = append(users, user)
			}
		} else {
			fmt.Println(index, user)
			users = append(users, user)
		}
	}

	return users

}

//
func ListDashboards(awsAccountID string, client *quicksight.QuickSight) ([]*quicksight.DashboardSummary, error) {
	input := &quicksight.ListDashboardsInput{
		AwsAccountId: aws.String(awsAccountID),
	}
	output, err := client.ListDashboards(input)
	if err != nil {
		return nil, err
	}

	return output.DashboardSummaryList, nil

}

func GetEmbedUrl(awsAccountID, region, roleName, dashboardID, userEmail string) *string {

	// Step 3: Get the embeddedURL

	client := InitClient(&region)

	namespace := "default"
	// Need to create separate client since dashboard region could be different from us-east-1 which is the user region
	userARN := "arn:aws:quicksight:" + region + ":" + awsAccountID + ":user/" + namespace + "/" + roleName + "/" + userEmail

	dashboardIdentityType := "QUICKSIGHT"

	eURLInput := quicksight.GetDashboardEmbedUrlInput{
		AwsAccountId: &awsAccountID,
		DashboardId:  &dashboardID,
		IdentityType: &dashboardIdentityType, //Needs to be QUICKSIGHT here and not IAM even  though an IAM role is being used that assumes the role
		UserArn:      &userARN,
	}

	eURLOutput, errEmbed := client.GetDashboardEmbedUrl(&eURLInput)

	if errEmbed != nil {
		fmt.Println("\nStep 3.2 - ", errEmbed.Error())
		return nil
	} else {
		return eURLOutput.EmbedUrl
	}
}
