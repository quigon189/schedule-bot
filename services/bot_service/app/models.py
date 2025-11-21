from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel


class CreateUserRequest(BaseModel):
    """Модель запроса создания пользователя в auth-service"""
    telegram_id: int  # обязательное поле
    username: Optional[str] = None  # может быть пустым
    full_name: str  # обязательное поле содержащее имя и фамилию


class AuthResponse(BaseModel):
    """
    Модель ответа auth-service
    дальше возможно заменим на 2 типо ответа:
        SuccessResponse для успешных ответов
        ErrorResponse для ответов с ошибкой
    """
    success: bool
    message: Optional[str] = None
    error: Optional[str] = None
    # содержит инфу о пользователе в успешных ответах
    data: Optional["UserResponse"] = None


class UserResponse(BaseModel):
    """Модель ответа от auth-service содержащая информацию о пользователе"""
    telegram_id: int
    username: Optional[str] = None
    full_name: str
    created_at: datetime
    updated_at: datetime
    is_active: bool
    roles: Optional[List["Role"]] = None  # содержит список ролей пользователя
    group: Optional[str] = None


class Role(BaseModel):
    """Модель роли в ответе от auth-service"""
    id: int
    name: str
    description: str


class TelegramUser(BaseModel):
    """Модель пользователя полученого из сообщения telegram"""
    id: int
    username: Optional[str] = None
    first_name: str
    last_name: Optional[str] = None

    @property
    def full_name(self) -> str:
        """получаем полное имя пользователя"""
        name = self.first_name
        if self.last_name:
            name += f" {self.last_name}"
        return name


class HealthCheck(BaseModel):
    status: str
    service: str
    timestampt: datetime
