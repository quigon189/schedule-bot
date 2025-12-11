from aiogram import F, Router
from aiogram.types import CallbackQuery

from app.models import TelegramUser
from app.services.user_service import user_service


register_callback_router = Router()


@register_callback_router.callback_query(F.data == "register_user")
async def register_user(callback_query: CallbackQuery):
    user = TelegramUser(
        id=callback_query.from_user.id,
        username=callback_query.from_user.username,
        first_name=callback_query.from_user.first_name,
        last_name=callback_query.from_user.last_name
    )

    reg_user = await user_service.register_user(user)
    if reg_user:
        await callback_query.answer("Спасибо за регистрацию!")
        message = f"""
✅ Добро пожаловать, {user.full_name}!

Вы успешно зарегистрированы в системе.

Теперь вы можете использовать все возможности бота!
"""

    else:
        await callback_query.answer("Ошибка регистрации")
        message = "При регистрации произошла ошибка, попробуйте повторить позднее."

    await callback_query.message.edit_text(message)
