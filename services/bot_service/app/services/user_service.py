from typing import Optional
from app.services.auth_service import auth_service
from app.models import UserResponse, TelegramUser


class UserService:
    """
    Класс для удобного взаимодействия с пользователями
    """
    async def register_user(self, telegram_user: TelegramUser) -> Optional[UserResponse]:
        """Регистрирует пользователя в системе"""
        user = await auth_service.create_user(
            telegram_id=telegram_user.id,
            username=telegram_user.username,
            full_name=telegram_user.full_name
        )
        return user

    async def get_user(self, user_id: int) -> Optional[UserResponse]:
        user = await auth_service.get_user(user_id)
        return user


user_service = UserService()
