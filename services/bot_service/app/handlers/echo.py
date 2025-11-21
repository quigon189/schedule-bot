from aiogram import Router
from aiogram.types import Message
from aiogram import F

from app.services.user_service import user_service

echo_router = Router()


@echo_router.message(F.text)
async def echo_handler(message: Message):
    user = await user_service.get_user(
        user_id=message.from_user.id
    )

    if user:
        response_text = f"""
📨 Вы написали: {message.text}

Вы зарегистрированны в системе
👤 Ваши данные:
Username: {user.username}
FullName: {user.full_name}
        """
    else:
        response_text = f"""
📨 Вы написали: {message.text}

Вы НЕ зарегистрированны в системе
        """

    await message.answer(response_text)
