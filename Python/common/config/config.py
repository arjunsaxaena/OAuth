from pydantic_settings import BaseSettings
import os

class Settings(BaseSettings):
    DB_URL: str
    MESSAGE_CENTRAL_AUTH_URL: str
    MESSAGE_CENTRAL_VALIDATE_URL: str
    MESSAGE_CENTRAL_SEND_OTP_URL: str
    MESSAGE_CENTRAL_CID: str
    MESSAGE_CENTRAL_KEY: str

    class Config:
        env_file = ".env"

settings = Settings() # Creating a settings() object which gets called everywhere we want to talk using the env variables
