from sqlalchemy.orm import sessionmaker
from Python.common.db.engine import engine

SessionLocal = sessionmaker (
    autocommit=False, # You have to manually commit transactions. Ensures changes aren’t accidentally saved.
    autoflush=False, # SQLAlchemy won’t automatically flush changes to the DB until you call commit().
    bind=engine # Tells SQLAlchemy which database connection to use. In our case the DB_URL in 'engine'.
) # Factory to create database sessions in SQLAlchemy (used to query and manipulate DB).

def get_db():
    db = SessionLocal()  # Create a session
    try:
        yield db        # Yield it (for FastAPI dependency injection)
    finally:
        db.close()      # Always close session to avoid DB connection leaks
                        # This pattern prevents connection leaks and ensures every request gets a fresh session.