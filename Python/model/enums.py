from enum import Enum

class Provider(str, Enum):
    PHONE = "phone"
    GOOGLE = "google"