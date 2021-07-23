package xsolla

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"testing"
)

func newTestClient() *Client {
	merchantID, _ := strconv.Atoi(os.Getenv("MerchantID"))
	projectID, _ := strconv.Atoi(os.Getenv("ProjectID"))
	return &Client {
		MerchantId:     merchantID,
		MerchantSecret: os.Getenv("MerchantSecret"),
		ProjectId:      projectID,
		ProjectSecret:  os.Getenv("ProjectSecret"),
		Sandbox :       true,
		Timeout :       100000000000,
	}
}

func newTestToken() *Token {
	UISettings := NewUISettings()
	var utm M
	var attr M
	userData := NewUserData("user_2",
		"john.smith@mail.com",
		"",
		utm,
		attr,)
	return &Token {
		User : userData,
		Settings : TokenSettings{
			ProjectId: 122221,
			UI : UISettings,
			Mode : "",
			ReturnURL : "",
		},

	}
}

func newTestAttribute() M{
	return M {
		"key" : "13",
		"list_of_values" : M {},
		"name" : M {
			"en": "rating",
		},
		"type" : "int",
	}
}

func TestClient_GetSubscriptionUserId(t *testing.T) {
	testClient := newTestClient()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	id := os.Getenv("testSubscriptionUserID")
	userID, err := testClient.GetSubscriptionUserId(id)
	if err != nil {
		t.Fatal(err)
	}
	if userID == "" {
		t.Fatal("Returned userID was nill")
	}
}


// No entry to test with ProjectID, will have to use 29054
// Additionally, part of the payload returns a date however there
// Seems to be an error when using json.Unmarshal() to parse the date
func TestClient_GetSubscription(t *testing.T) {
	testClient := newTestClient()
	projectID, _ := strconv.Atoi(os.Getenv("TestClientProjectID"))
	testClient.ProjectId = projectID
	var resPayLoad *Subscription
	var err error
	subscriptionID, _ := strconv.Atoi(os.Getenv("SubscriptionID"))
	resPayLoad, err = testClient.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}
	if resPayLoad == nil {
		t.Fatal("Subscription was nil")
	}
}

func TestClient_GetUser(t *testing.T) {
	testClient := newTestClient()
	var user *User
	var err error
	user, err = testClient.GetUser(os.Getenv("TestUserID"))
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("User is nill")
	}
}
////
func TestClient_GetTransaction(t *testing.T) {
	testClient := newTestClient()
	id := os.Getenv("TestClientID")
	transaction, err := testClient.GetTransaction(id)
	if err != nil {
		t.Fatal(err)
	}
	if transaction == nil {
		t.Fatal("Transaction is nil")
	}
}

func TestClient_CreateToken(t *testing.T) {
	testClient := newTestClient()
	testToken := newTestToken()
	token, err := testClient.CreateToken(testToken)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("Token is empty")
	}
}

// According to the xsolla documentation, this will
// API endpoint will be removed https://developers.xsolla.com/store-api/v1/attributes/user-attributes/create-attribute
func TestClient_CreateUserAttribute(t *testing.T) {
	testClient := newTestClient()
	testAttribute := newTestAttribute()
	id, err := testClient.CreateUserAttribute(testAttribute)
	if err != nil {
		t.Fatal(err)
	}
	if id == 0 {
		t.Fatal("User attribute failed to be created")
	}
}