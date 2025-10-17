from pydantic import BaseModel

class OtpRequest(BaseModel):
    phone: str
    country: int = 91
    otp_length: int = 6
    flow_type: str = "SMS"

class VerifyOtp(BaseModel):
    phone: str
    otp: str
    verification_id: str