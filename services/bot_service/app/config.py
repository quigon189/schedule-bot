from pydantic_settings import BaseSettings
from typing import Optional


class Settings(BaseSettings):
    TELEGRAM_BOT_TOKEN: str = ""

    DEBUG: bool = False

    AUTH_SERVICE_URL: str = "http://auth-service:8080"
    AUTH_SERVICE_TIMEOUT: int = 30

    SCHEDULES_SERVICE_URL: Optional[str] = None

    class Config:
        env_file = ".env"
        case_sensetive = True


def get_settings() -> Settings:
    return Settings()


settings = get_settings()
