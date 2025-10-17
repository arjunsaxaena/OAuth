import uuid # For uuid generation
from sqlalchemy.sql import func # For timestamp
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy import Column, String, Enum, JSON, DateTime # SQLAlchemy column types.
from pydantic import BaseModel
from Python.common.db.base import Base
from Python.common.db.tables import Tables
from Python.model.enums import Provider

class User(Base):
    __tablename__ = Tables.USERS
    id = Column(UUID, primary_key=True, default=uuid.uuid4)
    name = Column(String, unique=False, nullable=True)
    phone = Column(String, unique=True, nullable=True)
    email = Column(String, unique=True, nullable=True)
    provider = Column(Enum(Provider), nullable=False)
    profile_picture = Column(String, nullable=True)
    meta = Column(JSON, nullable=True)
    created_at = Column(DateTime, nullable=False, server_default=func.now())
    updated_at = Column(DateTime, nullable=False, server_default=func.now(), onupdate=func.now())

class LoginRequest(BaseModel):
    phone: str
    otp: str
    verification_id: str

class UserUpdate(BaseModel):
    name: str
    phone: str
    email: str
    profile_picture: str
    meta: str
