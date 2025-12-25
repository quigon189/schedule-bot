from aiogram import Router, types
from aiogram import F
from app.models import UserResponse
from app.keyboards.user_keyboards import get_main_menu_keyboard, get_schedule_menu_keyboard, get_ticket_menu_keyboard
from app.services.schedule_service import schedule_service
from datetime import datetime
import logging

user_router = Router()
logger = logging.getLogger(__name__)

@user_router.message(F.text == "📋 Профиль")
async def profile_button(message: types.Message, user: UserResponse):
    # Формируем текст профиля
    roles_text = ", ".join(user.roles_list) if user.roles_list else "нет ролей"
    
    profile_text = (
        f"👤 Ваш профиль\n\n"
        f"📝 Имя: {user.full_name}\n"
        f"🆔 ID: {user.telegram_id}\n"
        f"📧 Username: @{user.username if user.username else 'нет'}\n"
        f"🎭 Роли: {roles_text}\n"
        f"📊 Статус: {'✅ Активен' if user.is_active else '❌ Неактивен'}\n"
    )
    
    if user.group:
        profile_text += f"📚 Группа: {user.group}\n"
    
    profile_text += f"📅 Регистрация: {user.created_at.strftime('%d.%m.%Y %H:%M')}"
    
    await message.answer(profile_text, reply_markup=get_main_menu_keyboard())

# TODO: Сделать обработку расписания (его вызова)
@user_router.message(F.text == "")
async def schedule_button(message: types.Message, user: UserResponse):
    # Формируем текст профиля
    roles_text = ", ".join(user.roles_list) if user.roles_list else "нет ролей"
    
    profile_text = (
        f"👤 Ваш профиль\n\n"
        f"📝 Имя: {user.full_name}\n"
        f"🆔 ID: {user.telegram_id}\n"
        f"📧 Username: @{user.username if user.username else 'нет'}\n"
        f"🎭 Роли: {roles_text}\n"
        f"📊 Статус: {'✅ Активен' if user.is_active else '❌ Неактивен'}\n"
    )
    
    if user.group:
        profile_text += f"📚 Группа: {user.group}\n"
    
    profile_text += f"📅 Регистрация: {user.created_at.strftime('%d.%m.%Y %H:%M')}"
    
    await message.answer(profile_text, reply_markup=get_main_menu_keyboard())
