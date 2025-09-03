import requests
import base64
from fastapi import APIRouter
from Python.common.config.config import settings
from Python.model.otp import OtpRequest, VerifyOtp

router = APIRouter(prefix = "/sso")

@router.post("/request-otp")
def request_otp(request: OtpRequest):
    phone = request.phone
    country_code = request.country_code
    otp_length = request.otp_length
    flow_type = request.flow_type
    auth_url = getattr(settings, "MESSAGE_CENTRAL_AUTH_URL", None)
    key = getattr(settings, "MESSAGE_CENTRAL_KEY", None)

    if key is None:
        raise HTTPException("Key not set")

    key_b64 = base64.b64encode(key.encode()).decode()

    auth_params = {
        "customer_id": getattr(settings, "MESSAGE_CENTRAL_CID", None),
        "key": key_b64,
        "scope": "NEW",
        "country_code": country_code,
        "email": getattr(settings, "MESSAGE_CENTRAL_EMAIL", None),
    }

    if not auth_params["customer_id"]:
        raise HTTPException("Customer ID not set")

    try:
        auth_response = requests.get(auth_url, params=auth_params, headers = {"accept": "*/*"}, timeout = 10)
        auth_response.raise_for_status()
        auth_token = auth_response.json().get("token")

        if not auth_token:
            raise HTTPException("Auth token not found from message central")

    except Exception as e:
        raise InternalServerError(str(e))

    send_otp_url = getattr(settings, "MESSAGE_CENTRAL_SEND_OTP_URL", None)
    if not send_otp_url:
        raise HTTPException("Send OTP URL not set")

    send_otp_params = {
        "countryCode": country_code,
        "flowType": flow_type,
        "otpLength": otp_length,
        "mobileNumber": phone,
    }

    try:
        send_otp_response = requests.post(
            send_otp_url,
            params=send_otp_params,
            headers = {"authToken": auth_token},
            timeout = 10,
        )
        send_otp_response.raise_for_status()
        data = send_otp_response.json()
        if data.get("responseCode") == 200:
            return {
                "message": "OTP sent successfully",
                "verification_id": data["data"]["verificationId"],
            }

        else:
            raise HTTPException("OTP sent failed", data.get("message", "OTP sent failed"))

    except Exception as e:
        error_details = None
        if 'send_otp_response' in locals():
            try:
                error_details = send_otp_response.json()
            except Exception:
                error_details = send_otp_response.text
        else:
            error_details = str(e)
        raise InternalServerError("OTP sent failed", error_details)