from sqlalchemy.orm import Session
from fastapi import APIRouter, HTTPException, Header, Depends
from Python.api.otp import verify_otp
from Python.common.middleware.jwt import generate_token
from Python.model.user import User, LoginRequest
from Python.model.filters import GetUserFilters # Filter to get user by id, name, email, phone, provider
from Python.model.enums import Provider
from Python.common.db.session import get_db
from Python.common.repository.repository import UserRepository

router = APIRouter(prefix="/sso", tags=["SSO"])

user_repo = UserRepository()

@router.post("/login")
def login(request: LoginRequest, db: Session = Depends(get_db), authorization: str = Header(...)):
    if not authorization.startswith("Bearer "):
        raise HTTPException(status_code=400, detail="invalid auth header")

    auth_token = authorization.split(" ", 1)[1]
    
    if not verify_otp(request.verification_id, request.otp, auth_token):
        raise HTTPException(status_code=400, detail="Invalid OTP")

    user_filter = GetUserFilters(phone=request.phone)

    users = user_repo.get(db=db, filters=user_filter)
    existing_user = users[0] if users else None

    if existing_user:
        token = generate_token(str(existing_user.id))

        return {
            "token": token,
            "message": "login successful"
        }

    new_user = user_repo.create(db=db, obj_in=User(phone=request.phone, provider=Provider.PHONE))

    token = generate_token(str(new_user.id))

    return {
        "token": token,
        "message": "registration successful"
    }

