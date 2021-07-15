package xsolla

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrInvalidSignature = errors.New("invalid webhook signature ")
)

const (
	NotificationTypePayment            = "payment"
	NotificationTypeUserValidation     = "user_validation"
	NotificationTypeCreateSubscription = "create_subscription"
	NotificationTypeUpdateSubscription = "update_subscription"
	NotificationTypeCancelSubscription = "cancel_subscription"
)

type Webhook struct {
	// Raw is the pure string response of the webhook parsed from the HTTP request.
	Raw string `json:"-"`

	CustomParams     M            `json:"custom_parameters"`
	NotificationType string       `json:"notification_type"`
	PaymentDetails   M            `json:"payment_details"`
	Purchase         M            `json:"purchase"`
	Refund           M            `json:"refund_details"`
	Subscription     Subscription `json:"subscription"`
	Transaction      M            `json:"transaction"`
	User             User         `json:"user"`
}

func ParseWebhook(req *http.Request, projectSecret string) (*Webhook, error) {
	signature := strings.Replace(req.Header.Get("Authorization"), "Signature ", "", 1)
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	hash := sha1.New()
	hash.Write(buf)
	hash.Write([]byte(projectSecret))
	sum := fmt.Sprintf("%x", hash.Sum(nil))
	if signature != sum {
		return nil, ErrInvalidSignature
	}
	var hook Webhook
	if err := json.Unmarshal(buf, &hook); err != nil {
		return nil, err
	}
	hook.Raw = string(buf)
	return &hook, nil
}
