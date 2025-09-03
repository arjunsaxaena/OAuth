import requests
import base64
import logging
from fastapi import APIRouter, HTTPException, status
from Python.common.config.config import settings
from Python.model.otp import OtpRequest, VerifyOtp

logging.basicConfig(level=logging.DEBUG)
logger = logging.getLogger(__name__)

router = APIRouter(prefix="/sso")

@router.post("/request-otp")
def request_otp(request: OtpRequest):
    phone = request.phone
    country = request.country
    otp_length = request.otp_length
    flow_type = request.flow_type
    auth_url = getattr(settings, "MESSAGE_CENTRAL_AUTH_URL", None)
    key = getattr(settings, "MESSAGE_CENTRAL_KEY", None)

    logger.debug(f"Request received: phone={phone}, country_code={country}, otp_length={otp_length}, flow_type={flow_type}")

    if key is None:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Key not set")

    key_b64 = base64.b64encode(key.encode()).decode()
    logger.debug(f"Encoded key: {key_b64}")

    auth_params = {
        "customerId": getattr(settings, "MESSAGE_CENTRAL_CID", None),
        "key": key_b64,
        "scope": "NEW",
        "country": country,
        "email": getattr(settings, "MESSAGE_CENTRAL_EMAIL", None),
    }

    logger.debug(f"Auth params: {auth_params}")

    if not auth_params["customerId"]:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Customer ID not set")

    try:
        logger.debug(f"Sending auth request to {auth_url}")
        auth_response = requests.get(auth_url, params=auth_params, headers={"accept": "*/*"}, timeout=10)
        logger.debug(f"Auth response status: {auth_response.status_code}")
        auth_response.raise_for_status()
        auth_token = auth_response.json().get("token")
        logger.debug(f"Auth token: {auth_token}")

        if not auth_token:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Auth token not found from message central")

    except Exception as e:
        logger.exception("Auth request failed")
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=str(e))

    send_otp_url = getattr(settings, "MESSAGE_CENTRAL_SEND_OTP_URL", None)
    if not send_otp_url:
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail="Send OTP URL not set")

    send_otp_params = {
        "countryCode": country,
        "flowType": flow_type,
        "otpLength": otp_length,
        "mobileNumber": phone,
    }

    logger.debug(f"Sending OTP request to {send_otp_url} with params: {send_otp_params}")

    try:
        send_otp_response = requests.post(
            send_otp_url,
            params=send_otp_params,
            headers={"authToken": auth_token},
            timeout=10,
        )
        logger.debug(f"Send OTP response status: {send_otp_response.status_code}")
        send_otp_response.raise_for_status()
        data = send_otp_response.json()
        logger.debug(f"Send OTP response data: {data}")

        if data.get("responseCode") == 200:
            return {
                "message": "OTP sent successfully",
                "verification_id": data["data"]["verificationId"],
            }
        else:
            raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=data.get("message", "OTP sent failed"))

    except Exception as e:
        error_details = None
        if 'send_otp_response' in locals():
            try:
                error_details = send_otp_response.json()
            except Exception:
                error_details = send_otp_response.text
        else:
            error_details = str(e)
        logger.exception(f"OTP sending failed: {error_details}")
        raise HTTPException(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=f"OTP sent failed: {error_details}")
