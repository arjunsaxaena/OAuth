import jwt
from typing import Dict, Optional
from datetime import datetime, timedelta
from Python.common.config.config import settings

def generate_token(user_id: str) -> str:
    payload = {
            "user_id": user_id,
            "exp": datetime.utcnow() + timedelta(hours=100),
            "iat": datetime.utcnow()
    }

    return jwt.encode(payload, settings.JWT_SECRET, algorithm="HS256")

def validate_token(token: str) -> Optional[Dict]:
    try:
        return jwt.decode(token, settings.JWT_SECRET, algorithm="HS256")
    except (jwt.ExpiredSignatureError, jwt.InvalidTokenError):
        return None
