package xsolla

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	ErrInvalidSignature = errors.New("invalid webhook signature ")
)

const (
	// https://developers.xsolla.com/api/v2/getting-started/#api_webhooks_webhooks_list
	NotificationTypeAFSBlackList           = "afs_black_list"
	NotificationTypeAFSReject              = "afs_reject"
	NotificationTypeCancelSubscription     = "cancel_subscription"
	NotificationTypeCreateSubscription     = "create_subscription"
	NotificationTypeGetPinCode             = "get_pincode"
	NotificationTypeNonRenewalSubscription = "non_renewal_subscription"
	NotificationTypePayment                = "payment"
	NotificationTypePaymentAccountAdd      = "payment_account_add"
	NotificationTypePaymentAccountRemove   = "payment_account_remove"
	NotificationTypeRedeemKey              = "redeem_key"
	NotificationTypeRefund                 = "refund"
	NotificationTypeUpdateSubscription     = "update_subscription"
	NotificationTypeUpgradeRefund          = "upgrade_refund"
	NotificationTypeUserBalanceOperation   = "user_balance_operation"
	NotificationTypeUserSearch             = "user_search"
	NotificationTypeUserValidation         = "user_validation"

	// https://developers.xsolla.com/api/v2/getting-started/#api_webhooks_errors
	WebhookErrCodeInvalidUser      = "INVALID_USER"
	WebhookErrCodeInvalidParameter = "INVALID_PARAMETER"
	WebhookErrCodeInvalidSignature = "INVALID_SIGNATURE"
	WebhookErrCodeIncorrectAmount  = "INCORRECT_AMOUNT"
	WebhookErrCodeIncorrectInvoice = "INCORRECT_INVOICE"
)

type Webhook struct {
	// Raw is the pure string response of the webhook parsed from the HTTP request.
	Raw string `json:"-"`

	CustomParams     M                   `json:"custom_parameters"`
	NotificationType string              `json:"notification_type"`
	PaymentDetails   M                   `json:"payment_details"`
	Purchase         M                   `json:"purchase"`
	Refund           M                   `json:"refund_details"`
	Subscription     WebhookSubscription `json:"subscription"`
	Transaction      M                   `json:"transaction"`
	User             User                `json:"user"`
}

// The webhook subscription structure differs than the regular API structure, for who knows what reason.
// https://developers.xsolla.com/api/v2/getting-started/#api_webhooks_created_subscription
type WebhookSubscription struct {
	ChargeAmount   float64   `json:"charge_amount"`
	Comment        string    `json:"comment"`
	Currency       string    `json:"currency"`
	DateCreate     time.Time `json:"date_create"`
	DateEnd        time.Time `json:"date_end"`
	DateNextCharge time.Time `json:"date_next_charge"`
	SubscriptionId int       `json:"subscription_id"`
	PlanId         string    `json:"plan_id"`
	ProductId      string    `json:"product_id"`
	Status         string    `json:"status"`
	Tags           []string  `json:"tags"`
	Trial          Trial     `json:"trial"`
	User           User      `json:"user"`
}

type WebhookError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// https://developers.xsolla.com/api/v2/getting-started/#api_webhooks_signing_requests
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
