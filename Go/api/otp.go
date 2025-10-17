package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"oauth/Go/common/logger"
	"oauth/Go/common/util"
	"oauth/Go/model"
)

func (h *Handler) RequestOTP(w http.ResponseWriter, r *http.Request) {
	var request model.OTPRequest
	
	if err := util.ReadJsonAndValidate(w, r, &request); err != nil {
		logger.Error("Failed to read/validate request: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	if h.config.MessageCentral.Key == "" {
		logger.Error("Message Central key not configured")
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Key not set")
		return
	}

	if h.config.MessageCentral.CID == "" {
		logger.Error("Message Central CID not configured")
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Customer ID not set")
		return
	}

	keyB64 := base64.StdEncoding.EncodeToString([]byte(h.config.MessageCentral.Key))

	authToken, err := h.getAuthToken(keyB64, request.Country)
	if err != nil {
		logger.Error("Failed to get auth token: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Auth request failed: %v", err))
		return
	}

	verificationID, err := h.sendOTP(authToken, request)
	if err != nil {
		logger.Error("Failed to send OTP: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("OTP sending failed: %v", err))
		return
	}

	response := map[string]interface{}{
		"message":        "OTP sent successfully",
		"verification_id": verificationID,
		"auth_token":     authToken,
	}

	util.WriteSuccessResponse(w, response)
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var request model.VerifyOTPRequest
	
	if err := util.ReadJsonAndValidate(w, r, &request); err != nil {
		logger.Error("Failed to read/validate verification request: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	isValid, err := h.verifyOTPCode(request.VerificationID, request.OTP, request.AuthToken)
	if err != nil {
		logger.Error("OTP verification failed: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("OTP verification failed: %v", err))
		return
	}

	response := map[string]interface{}{
		"valid": isValid,
		"message": func() string {
			if isValid {
				return "OTP verified successfully"
			}
			return "Invalid OTP"
		}(),
	}

	util.WriteSuccessResponse(w, response)
}

func (h *Handler) getAuthToken(keyB64 string, country int) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequest("GET", h.config.MessageCentral.AuthURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	q := req.URL.Query()
	q.Add("customerId", h.config.MessageCentral.CID)
	q.Add("key", keyB64)
	q.Add("scope", "NEW")
	q.Add("country", strconv.Itoa(country))
	req.URL.RawQuery = q.Encode()

	req.Header.Set("accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResponse struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	if authResponse.Token == "" {
		return "", fmt.Errorf("auth token not found from message central")
	}

	return authResponse.Token, nil
}

func (h *Handler) sendOTP(authToken string, request model.OTPRequest) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequest("POST", h.config.MessageCentral.SendOTPURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create send OTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("countryCode", strconv.Itoa(request.Country))
	q.Add("flowType", request.FlowType)
	q.Add("otpLength", strconv.Itoa(request.OTPLength))
	q.Add("mobileNumber", request.Phone)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("authToken", authToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send OTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read send OTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("send OTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var otpResponse struct {
		ResponseCode int `json:"responseCode"`
		Message      string `json:"message"`
		Data         struct {
			VerificationID string `json:"verificationId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &otpResponse); err != nil {
		return "", fmt.Errorf("failed to decode send OTP response: %w", err)
	}

	if otpResponse.ResponseCode != 200 {
		return "", fmt.Errorf("OTP send failed: %s", otpResponse.Message)
	}

	return otpResponse.Data.VerificationID, nil
}

func (h *Handler) verifyOTPCode(verificationID, otp, authToken string) (bool, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	
	req, err := http.NewRequest("GET", h.config.MessageCentral.ValidateURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create verify OTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("verificationId", verificationID)
	q.Add("code", otp)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("authToken", authToken)

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("verify OTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read verify OTP response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Error("OTP validation failed with status %d: %s", resp.StatusCode, string(body))
		return false, nil
	}

	var verifyResponse struct {
		ResponseCode int `json:"responseCode"`
	}

	if err := json.Unmarshal(body, &verifyResponse); err != nil {
		return false, fmt.Errorf("failed to decode verify OTP response: %w", err)
	}

	return verifyResponse.ResponseCode == 200, nil
}