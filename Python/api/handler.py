from fastapi import APIRouter
from Python.api.otp import router as otp_router
from Python.api.user import router as user_router

router = APIRouter()

router.include_router(otp_router)
router.include_router(user_router)
