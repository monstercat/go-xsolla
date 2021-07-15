package xsolla

import (
	"time"
)

const (
	EndpointMerchant = "https://api.xsolla.com/merchant/v2/merchants"
	EndpointProject  = "https://api.xsolla.com/merchant/v2/projects"

	modeSandbox = "sandbox"

	SubscriptionActive      = "active"
	SubscriptionCanceled    = "canceled"
	SubscriptionEnded       = "ended"
	SubscriptionNonRenewing = "non_renewing"
)

// https://developers.xsolla.com/api/v2/getting-started/#api_errors_handling
type RequestError struct {
	Raw             string `json:"-"`
	Code            int    `json:"http_status_code"`
	Message         string `json:"message"`
	ExtendedMessage string `json:"extended_message"`
	RequestId       string `json:"request_id"`
}

func (e *RequestError) Error() string {
	if e.Message != "" {
		return e.Message
	} else if e.Raw != "" {
		return e.Raw
	}
	return "empty xsolla request error"
}

type M map[string]interface{}

type User struct {
	Country string `json:"country"`
	Email   string `json:"email"`
	Id      string `json:"id"`
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
}

type Subscription struct {
	ChargeAmount   float64   `json:"charge_amount"`
	Comment        string    `json:"comment"`
	Currency       string    `json:"currency"`
	DateCreate     time.Time `json:"date_create"`
	DateEnd        time.Time `json:"date_end"`
	DateNextCharge time.Time `json:"date_next_charge"`
	Id             int       `json:"id"`
	Plan           Plan      `json:"plan"`
	Status         string    `json:"status"`
	Tags           []string  `json:"tags"`
	Trial          Trial     `json:"trial"`
	User           User      `json:"user"`
	// Product (not sure what this is).
}

type Plan struct {
	ExternalId string `json:"external_id"`
	Id         int    `json:"id"`
}

type Trial struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type TokenSettings struct {
	ProjectId int                    `json:"project_id"`
	UI        map[string]interface{} `json:"ui"`
	Mode      string                 `json:"mode,omitempty"`
	ReturnURL string                 `json:"return_url,omitempty"`
}

type Token struct {
	User         M                 `json:"user"`
	CustomParams M                 `json:"custom_parameters,omitempty"`
	Settings     TokenSettings     `json:"settings"`
	Purchase     *PurchaseSettings `json:"purchase,omitempty"`
}

type PurchaseSettings struct {
	Subscription PurchaseSubscription `json:"subscription"`
	Description  PurchaseDescription  `json:"description"`
	CouponCode   *PurchaseCouponCode  `json:"coupon_code,omitempty"`
}

type PurchaseDescription struct {
	Value string `json:"value"`
}

type PurchaseCouponCode struct {
	Value string `json:"value"`
}

type PurchaseSubscription struct {
	PlanId    string `json:"plan_id,omitempty"`
	ProductId string `json:"product_id,omitempty"`

	// Subscription plans (array) to show in the payment UI.
	AvailablePlans []string `json:"available_plans,omitempty"`

	// The type of operation applied to the user’s subscription plan. To change the subscription plan, pass the
	// change_plan value. You need to specify the new plan ID in the purchase.subscription.plan_id parameter.
	Operation string `json:"operation,omitempty"`

	// Currency of the subscription plan to use in all calculations.
	Currency string `json:"currency,omitempty"`

	// Trial period in days.
	TrialDays int `json:"trial_days,omitempty"`
}

type TransactionSubscriptionDetails struct {
	IsPaymentFromSubscription bool `json:"is_payment_from_subscription"`
	IsSubscriptionCreated     bool `json:"is_subscription_created"`
}

type Transaction struct {
	CustomerDetails     map[string]interface{}         `json:"customer_details"`
	TransactionDetails  map[string]interface{}         `json:"transaction_details"`
	SubscriptionDetails TransactionSubscriptionDetails `json:"subscription_details"`
}

func NewUserData(id, email, promo string, utm, attr M) M {
	data := M{
		"id": M{
			"value":  id,
			"hidden": true,
		},
		"email": M{
			"value": email,
		},
		"country": M{
			"allow_modify": true,
		},
	}

	if promo != "" {
		attr["promo"] = promo
	}

	if len(attr) > 0 {
		data["attributes"] = attr
	}

	if len(utm) > 0 {
		data["utm"] = utm
	}

	return data
}

func NewCustomParams(active, reg time.Time, tf bool) M {
	return M{
		"registration_date":       reg,
		"active_date":             active,
		"additional_verification": tf,
	}
}

func NewUISettings() M {
	return M{
		"version": "desktop",
		"desktop": M{
			"header": M{
				"type":         "compact",
				"visible_name": true,
			},
			"subscription_list": M{},
		},
		"mobile": M{
			"footer": M{
				"is_visible": false,
			},
		},
		"components": M{
			"virtual_currency": M{
				"hidden": true,
			},
		},
	}
}

func NewUTM(source, campaign, term, content string) M {
	return M{
		"utm_source":   source,
		"utm_campaign": campaign,
		"utm_term":     term,
		"utm_content":  content,
	}
}
