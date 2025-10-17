from sqlalchemy import create_engine
from Python.common.config.config import settings

engine = create_engine(settings.DB_URL) # The engine is what actually talks to your DB (Postgres, MySQL, SQLite, etc.).