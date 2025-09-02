from pydantic import BaseModel
from uuid import UUID
from typing import Optional

class UserFilters(BaseModel):
    id: Optional[UUID] = None
    name: Optional[str] = None
    email: Optional[str] = None
    phone: Optional[str] = None
    provider: Optional[str] = None