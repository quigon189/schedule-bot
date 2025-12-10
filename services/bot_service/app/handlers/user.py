from aiogram import Router, types
from aiogram.filters import Text
from services.auth_service import AuthService
from keyboards.user_keyboards import get_schedule_menu_keyboard
import logging

router = Router()
auth_service = AuthService()

async def IsRegistred(telegram_id: int) -> bool:
    """Проверка регистрации пользователя"""
    user = await auth_service.get_user(telegram_id)
    return user is not None


@router.message(Text("📋 Профиль"))
async def profile_button(message: types.Message):
    """Обработчик кнопки Профиль"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    user = await auth_service.get_user(telegram_id)
    if user:
        await message.answer(
            f"📋 *Ваш профиль*\n\n"
            f"👤 *Имя:* {user.full_name}\n"
            f"🆔 *Telegram ID:* {user.telegram_id}\n"
            f"📧 *Username:* @{user.username if user.username else 'не указан'}\n"
            f"📅 *Дата регистрации:* {user.created_at}",
            parse_mode="Markdown"
        )


@router.message(Text("📅 Расписание"))
async def schedule_button(message: types.Message):
    """Обработчик кнопки Расписание"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await message.answer(
        "📅 *Расписание*\n\n"
        "Выберите период:",
        parse_mode="Markdown",
        reply_markup=get_schedule_menu_keyboard()
    )


@router.message(Text("⚙️ Настройки"))
async def settings_button(message: types.Message):
    """Обработчик кнопки Настройки"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await message.answer(
        "⚙️ *Настройки*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown"
    )