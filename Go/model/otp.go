package model

type OTPRequest struct {
	Phone string `json:"phone"`
	Country int `json:"country" default:"91"`
	OTPLength int `json:"otp_length" default:"6"`
	FlowType string `json:"flow_type" default:"SMS"`
}

type OTPResponse struct {
	Phone string `json:"phone"`
	OTP string `json:"otp"`
	VerificationID string `json:"verification_id"`
}

type VerifyOTPRequest struct {
	Phone string `json:"phone"`
	OTP string `json:"otp"`
	VerificationID string `json:"verification_id"`
	AuthToken string `json:"auth_token"`
}