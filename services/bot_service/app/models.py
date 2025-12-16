from datetime import datetime
from typing import Any, List, Literal, Optional
from pydantic import BaseModel, Field

Roles = Literal["student", "teacher", "manager"]


class CreateUserRequest(BaseModel):
    """Модель запроса создания пользователя в auth-service"""
    code: str
    telegram_id: int  # обязательное поле
    username: Optional[str] = None  # может быть пустым
    full_name: str  # обязательное поле содержащее имя и фамилию


class CreateCodeRequest(BaseModel):
    role_name: Roles
    group_name: Optional[str]
    created_by: int
    max_uses: Optional[int]
    expiration: Optional[int]


class ServiceResponse(BaseModel):
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
    data: Any


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

    @property
    def roles_list(self) -> list:
        return [role.name for role in self.roles] if self.roles else []


class CodeResponse(BaseModel):
    id: int
    code: str
    role: "Role"
    group_name: Optional[str] = None
    max_uses: int
    created_by: Optional[UserResponse]
    expires_at: datetime
    created_at: datetime


class Role(BaseModel):
    """Модель роли в ответе от auth-service"""
    id: int
    name: str
    description: str


class GroupScheduleRequest(BaseModel):
    """Модель запроса расписания"""
    academic_year: str = Field(pattern=r"^(\d{4})/(\d{4})$")
    half_year: int = Field(ge=1, le=2)
    group_name: str


class GroupScheduleResponse(BaseModel):
    """Ответ от сервиса расписаний"""
    academic_year: str = Field(pattern=r"^(\d{4})/(\d{4})$")
    half_year: int = Field(ge=1, le=2)
    group_name: str
    semester: int = Field(ge=1, le=10)
    schedule_img_url: str
    created_at: datetime


class ScheduleChangesResponse(BaseModel):
    """Модель изменений расписания"""
    id: int
    date: datetime
    description: str
    image_urls: List[str]
    created_at: datetime


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
