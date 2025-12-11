from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message

from app.middlewaries.admin import AdminMiddleware
from app.models import UserResponse
from app.services.auth_service import auth_service


admin_router = Router()
admin_router.message.middleware(AdminMiddleware())


@admin_router.message(Command("code"))
async def create_code(message: Message, user: UserResponse):
    code = await auth_service.create_registration_code(
        role="teacher",
        created_by=user.telegram_id
    )
    if code:
        result_text = f"""
Код доступа для преподавателя создан
Код: {code.code}
Количество использований: {code.max_uses}
Действителен до: {code.expires_at.strftime("%Y-%y-%d %H:%M:%S")}
        """
    else:
        result_text = "При создании кода возникла ошибка"

    await message.answer(result_text)
