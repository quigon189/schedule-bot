from aiogram import Router, types
from aiogram.filters import Command
from aiogram.types import Message
from services.auth_service import AuthService
from services.schedule_service import ScheduleService
from keyboards.user_keyboards import get_main_menu_keyboard
from web.bot import bot
import logging

router = Router()
auth_service = AuthService()
schedule_service = ScheduleService()

async def IsRegistred(telegram_id: int) -> bool:
    """Проверка регистрации пользователя"""
    user = await auth_service.get_user(telegram_id)
    return user is not None


@router.message(Command("start"))
async def cmd_start(message: Message):
    """Обработчик команды /start"""
    telegram_id = message.from_user.id
    username = message.from_user.username
    full_name = message.from_user.full_name
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        # Если не зарегистрирован, создаем пользователя
        user = await auth_service.create_user(telegram_id, username, full_name)
        if user:
            await message.answer(
                f"Привет, {full_name}! Вы успешно зарегистрированы.\n"
                f"Используйте команду /menu для доступа к функциям.",
                reply_markup=get_main_menu_keyboard()
            )
        else:
            await message.answer("Ошибка регистрации. Пожалуйста, попробуйте позже.")
    else:
        await message.answer(
            f"С возвращением, {full_name}!",
            reply_markup=get_main_menu_keyboard()
        )


@router.message(Command("menu"))
async def cmd_menu(message: Message):
    """Обработчик команды /menu"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await message.answer(
        "Главное меню:",
        reply_markup=get_main_menu_keyboard()
    )


@router.message(Command("profile"))
async def cmd_profile(message: Message):
    """Обработчик команды /profile"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    # Получаем данные пользователя
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
    else:
        await message.answer("Не удалось загрузить данные профиля.")


@router.message(Command("schedule"))
async def cmd_schedule(message: Message):
    """Обработчик команды /schedule"""
    telegram_id = message.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    # Здесь можно добавить логику для получения группы пользователя
    # и запроса расписания
    await message.answer(
        "📅 *Расписание*\n\n"
        "Выберите действие:",
        parse_mode="Markdown",
        reply_markup=get_main_menu_keyboard()
    )