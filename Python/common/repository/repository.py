from Python.common.repository.base import BaseRepository
from Python.model.user import User

class UserRepository(BaseRepository):
    def __init__(self):
        super().__init__(User)