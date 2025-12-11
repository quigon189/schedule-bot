from typing import Optional
from aiogram import BaseMiddleware
from aiogram.types import Message, CallbackQuery

from app.models import UserResponse
from app.services.user_service import user_service


class AdminMiddleware(BaseMiddleware):
    async def __call__(self, handler, event: Message | CallbackQuery, data: dict):
        user: Optional[UserResponse] = data.get('user')
        if not user:
            user = await user_service.get_user(event.from_user.id)

        if not user:
            await event.answer("Ошибка при обработке запроса")
            return

        if "admin" in user.roles_list:
            return await handler(event, data)
        else:
            await event.answer("У вас недостаточно прав")
            return
