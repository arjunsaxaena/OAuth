from pydantic import BaseModel

class OtpRequest(BaseModel):
    phone: str
    country_code: str = 91
    otp_length: int = 6
    flow_type: str = "SMS"

class VerifyOtp(BaseModel):
    phone: str
    otp: str
    verification_id: str