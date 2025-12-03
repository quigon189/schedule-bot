from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    TELEGRAM_BOT_TOKEN: str = ""

    DEBUG: bool = False

    AUTH_SERVICE_URL: str = "http://auth-service:8080"
    AUTH_SERVICE_TIMEOUT: int = 30

    SCHEDULE_SERVICE_URL: str = "http://schedule-service:8080"
    SCHEDULE_SERVICE_TIMEOUT: int = 30

    class Config:
        env_file = ".env"
        case_sensetive = True


def get_settings() -> Settings:
    return Settings()


settings = get_settings()
