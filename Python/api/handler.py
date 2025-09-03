from fastapi import APIRouter
from Python.api.otp import router as otp_router

router = APIRouter()

router.include_router(otp_router)