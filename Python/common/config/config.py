from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    DB_URL: str = "postgresql://postgres:secret@localhost:5432/oauth?sslmode=disable"

    class Config:
        env_file = "../.env"

settings = Settings() # Creating a settings() object which gets called everywhere we want to talk using the env variables