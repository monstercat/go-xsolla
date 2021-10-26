package xsolla

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func newTestClient() *Client {
	merchantId, _ := strconv.Atoi(os.Getenv("MerchantId"))
	projectId, _ := strconv.Atoi(os.Getenv("ProjectId"))
	return &Client{
		MerchantId:     merchantId,
		MerchantSecret: os.Getenv("MerchantSecret"),
		ProjectId:      projectId,
		ProjectSecret:  os.Getenv("ProjectSecret"),
		Sandbox:        true,
		Timeout:        time.Second * 10,
	}
}

func newTestAttribute() M {
	return M{
		"key":            "13",
		"list_of_values": M{},
		"name": M{
			"en": "rating",
		},
		"type": "int",
	}
}

func TestClient_GetSubscriptionUserId(t *testing.T) {
	testClient := newTestClient()
	id := os.Getenv("testSubscriptionUserId")
	userId, err := testClient.GetSubscriptionUserId(id)
	if err != nil {
		t.Fatal(err)
	}
	if userId == "" {
		t.Fatal("Returned userId was nil")
	}
}

func TestClient_GetSubscription(t *testing.T) {
	projectId, _ := strconv.Atoi(os.Getenv("TestClientProjectId"))
	subscriptionId, _ := strconv.Atoi(os.Getenv("SubscriptionId"))
	testClient := newTestClient()
	testClient.ProjectId = projectId
	resPayLoad, err := testClient.GetSubscription(subscriptionId)
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
	user, err = testClient.GetUser(os.Getenv("TestUserId"))
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("User is nill")
	}
}

func TestClient_GetTransaction(t *testing.T) {
	testClient := newTestClient()
	id := os.Getenv("TestClientId")
	transaction, err := testClient.GetTransaction(id)
	if err != nil {
		t.Fatal(err)
	}
	if transaction == nil {
		t.Fatal("Transaction is nil")
	}
}

func TestClient_UpdateSubscription(t *testing.T) {
	subscriptionID, _ := strconv.Atoi(os.Getenv("SubscriptionId"))
	userID := os.Getenv("TestUserId")
	testClient := newTestClient()

	res, err := testClient.UpdateSubscription(userID, subscriptionID, &XsollaUpdatePayload{status: "canceled"})
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Subscription is nil, update subscription failed")
	}
	if res.Status != "canceled" {
		t.Fatal("Update subscription did not go through, please try again")
	}

}

func TestClient_CreateToken(t *testing.T) {
	var utm, attr M
	client := newTestClient()
	token := &Token{
		User:     NewUserData("user_2", "john.smith@mail.com", "", utm, attr),
		Settings: client.NewTokenSettings(),
	}
	str, err := client.CreateToken(token)
	if err != nil {
		t.Fatal(err)
	} else if str == "" {
		t.Fatal("Token response was empty")
	}
}
