package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

// Zarinpal is the base struct for zarinpal payment
// gateway, one shall not create or manipulate instances
// if this struct manually and just use provided methods
// to woek with it.
type Zarinpal struct {
	MerchantID      string
	Sandbox         bool
	APIEndpoint     string
	PaymentEndpoint string
}

type paymentRequestReqBody struct {
	MerchantID  string `json:"merchant_id"`
	Amount      int    `json:"amount"`
	CallbackURL string `json:"callback_url"`
	Description string `json:"description"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
}

type paymentRequestResp struct {
	Data struct {
		Status    int    `json:"code"`
		Authority string `json:"authority"`
	} `json:"data"`
}

type paymentVerificationReqBody struct {
	MerchantID string `json:"merchant_id"`
	Authority  string `json:"authority"`
	Amount     int    `json:"amount"`
}

type paymentVerificationResp struct {
	Data struct {
		StatusCode int    `json:"code"`
		Message    string `json:"message"`
		CardHash   string `json:"card_hash"`
		CardPan    string `json:"card_pan"`
		RefID      int    `json:"ref_id"`
		FeeType    string `json:"fee_type"`
		Fee        int    `json:"fee"`
	} `json:"data"`
	Errors []any
}

type unverifiedTransactionsReqBody struct {
	MerchantID string
}

// UnverifiedAuthority is the base struct for Authorities in unverifiedTransactionsResp
type UnverifiedAuthority struct {
	Authority   string
	Amount      int
	Channel     string
	CallbackURL string
	Referer     string
	Email       string
	CellPhone   string
	Date        string // ToDo Check type to be date
}

type unverifiedTransactionsResp struct {
	Status      int
	Authorities []UnverifiedAuthority
}

type refreshAuthorityReqBody struct {
	MerchantID string
	Authority  string
	ExpireIn   int
}

type refreshAuthorityResp struct {
	Status int
}

// NewZarinpal creates a new instance of zarinpal payment
// gateway with provided configs. It also tries to validate
// provided configs.
func NewZarinpal(merchantID string, sandbox bool) (*Zarinpal, error) {
	if len(merchantID) != 36 {
		return nil, errors.New("MerchantID must be 36 characters")
	}
	apiEndPoint := "https://payment.zarinpal.com/pg/v4/payment/"
	paymentEndpoint := "https://payment.zarinpal.com/pg/StartPay/"
	if sandbox {
		apiEndPoint = "https://sandbox.zarinpal.com/pg/v4/payment/"
		paymentEndpoint = "https://sandbox.zarinpal.com/pg/StartPay/"
	}
	return &Zarinpal{
		MerchantID:      merchantID,
		Sandbox:         sandbox,
		APIEndpoint:     apiEndPoint,
		PaymentEndpoint: paymentEndpoint,
	}, nil
}

// NewPaymentRequest gets a payment url from Zarinpal.
// amount is in Tomans (not Rials) format.
// email and mobile are optional.
//
// If error is not nil, you can check statusCode for
// specific error handling based on Zarinpal error codes.
// If statusCode is not 100, it means Zarinpal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinpal *Zarinpal) NewPaymentRequest(amount int, callbackURL, description, email, mobile string) (paymentURL, authority string, statusCode int, err error) {
	if amount < 1 {
		err = errors.New("amount must be a positive number")
		return
	}
	if callbackURL == "" {
		err = errors.New("callbackURL should not be empty")
		return
	}
	if description == "" {
		err = errors.New("description should not be empty")
		return
	}
	paymentRequest := paymentRequestReqBody{
		MerchantID:  zarinpal.MerchantID,
		Amount:      amount,
		CallbackURL: callbackURL,
		Description: description,
		Email:       email,
		Mobile:      mobile,
	}
	var resp paymentRequestResp
	err = zarinpal.request("request.json", &paymentRequest, &resp)
	if err != nil {
		return
	}
	statusCode = resp.Data.Status
	if resp.Data.Status == 100 {
		authority = resp.Data.Authority
		paymentURL = zarinpal.PaymentEndpoint + resp.Data.Authority
	} else {
		err = errors.New(strconv.Itoa(resp.Data.Status))
	}
	return
}

// PaymentVerification verifies if a payment was done successfully, Authority of the
// payment request should be passed to this method alongside its Amount in Tomans.
//
// If error is not nil, you can check statusCode for
// specific error handling based on Zarinpal error codes.
// If statusCode is not 100, it means Zarinpal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinpal *Zarinpal) PaymentVerification(amount int, authority string) (verified bool, refID string, statusCode int, err error) {
	if amount <= 0 {
		err = errors.New("amount must be a positive number")
		return
	}
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}
	paymentVerification := paymentVerificationReqBody{
		MerchantID: zarinpal.MerchantID,
		Amount:     amount,
		Authority:  authority,
	}
	var resp paymentVerificationResp
	err = zarinpal.request("verify.json", &paymentVerification, &resp)
	if err != nil {
		return
	}
	statusCode = resp.Data.StatusCode
	if resp.Data.StatusCode == 100 {
		verified = true
		refID = string(resp.Data.RefID)
	} else {
		err = errors.New(strconv.Itoa(resp.Data.StatusCode))
	}
	return
}

// UnverifiedTransactions gets unverified transactions.
//
// If error is not nil, you can check statusCode for
// specific error handling based on Zarinpal error codes.
// If statusCode is not 100, it means Zarinpal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinpal *Zarinpal) UnverifiedTransactions() (authorities []UnverifiedAuthority, statusCode int, err error) {
	unverifiedTransactions := unverifiedTransactionsReqBody{
		MerchantID: zarinpal.MerchantID,
	}

	var resp unverifiedTransactionsResp
	err = zarinpal.request("UnverifiedTransactions.json", &unverifiedTransactions, &resp)
	if err != nil {
		return
	}

	if resp.Status == 100 {
		statusCode = resp.Status
		authorities = resp.Authorities
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

// RefreshAuthority update authority expiration time.\n
// expire should be number between [1800,3888000] seconds.
//
// If error is not nil, you can check statusCode for
// specific error handling based on Zarinpal error codes.
// If statusCode is not 100, it means Zarinpal raised an error
// on their end and you can check the error code and its reason
// based on their documentation placed in
// https://github.com/ZarinPal-Lab/Documentation-PaymentGateway/archive/master.zip
func (zarinpal *Zarinpal) RefreshAuthority(authority string, expire int) (statusCode int, err error) {
	if authority == "" {
		err = errors.New("authority should not be empty")
		return
	}
	if expire < 1800 {
		err = errors.New("expire must be at least 1800")
		return
	} else if expire > 3888000 {
		err = errors.New("expire must not be greater than 3888000")
		return
	}

	refreshAuthority := refreshAuthorityReqBody{
		MerchantID: zarinpal.MerchantID,
		Authority:  authority,
		ExpireIn:   expire,
	}
	var resp refreshAuthorityResp
	err = zarinpal.request("RefreshAuthority.json", &refreshAuthority, &resp)
	if err != nil {
		return
	}
	if resp.Status == 100 {
		statusCode = resp.Status
	} else {
		err = errors.New(strconv.Itoa(resp.Status))
	}
	return
}

func (zarinpal *Zarinpal) request(method string, data interface{}, res interface{}) error {
	reqBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", zarinpal.APIEndpoint+method, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(body))
	err = json.Unmarshal(body, res)
	if err != nil {
		err = errors.New("zarinpal invalid json response")
		return err
	}
	return nil
}
