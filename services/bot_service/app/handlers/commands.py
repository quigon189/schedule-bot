from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message
from app.models import UserResponse
from app.keyboards.user_keyboards import get_main_menu_keyboard
import logging
from datetime import datetime
from app.services.schedule_service import schedule_service

logger = logging.getLogger(__name__)

router = Router()

@router.message(Command("start"))
async def cmd_start(message: Message, user: UserResponse = None):
    if user:
        await message.answer(
            f"✅ Привет, {user.full_name}!\n\n"
            "Вы уже зарегистрированы в системе.\n"
            "Используйте меню для навигации.",
            reply_markup=get_main_menu_keyboard()
        )
    else:
        await message.answer(
            "👋 Добро пожаловать!\n\n"
            "Для использования бота вам нужно зарегистрироваться.\n"
            "Получите код регистрации у администратора."
        )

@router.message(Command("menu"))
async def cmd_menu(message: Message, user: UserResponse = None):
    if not user:
        await message.answer("Сначала зарегистрируйтесь!")
        return
    
    await message.answer(
        "🏠 Главное меню\n\nВыберите действие:",
        reply_markup=get_main_menu_keyboard()
    )

@router.message(Command("profile"))
async def cmd_profile(message: Message, user: UserResponse = None):
    if not user:
        await message.answer("Сначала зарегистрируйтесь!")
        return
    
    roles_text = ", ".join(user.roles_list) if user.roles_list else "нет ролей"
    profile_text = (
        f"👤 Ваш профиль\n\n"
        f"Имя: {user.full_name}\n"
        f"ID: {user.telegram_id}\n"
        f"Username: @{user.username if user.username else 'нет'}\n"
        f"Роли: {roles_text}\n"
    )
    
    if user.group:
        profile_text += f"Группа: {user.group}\n"
    
    profile_text += f"Дата регистрации: {user.created_at}"
    
    await message.answer(profile_text, reply_markup=get_main_menu_keyboard())

@router.message(Command("help"))
async def cmd_help(message: Message, user: UserResponse = None):
    help_text = (
        "📚 Доступные команды:\n\n"
        "/start - Начать работу с ботом\n"
        "/menu - Главное меню\n"
        "/profile - Показать профиль\n"
        "/help - Эта справка\n\n"
    )
    
    if user and "admin" in user.roles_list:
        help_text += "⚙️ Команды администратора:\n"
        help_text += "/code - Создать код регистрации\n"
    
    await message.answer(help_text)

@router.message(Command("schedule"))
async def cmd_schedule(message: Message, user: UserResponse = None):
    if not user:
        await message.answer("Сначала зарегистрируйтесь!")
        return
    
    # Если студент с группой, сразу показываем его расписание
    if hasattr(user, 'roles_list') and 'student' in user.roles_list and user.group:
        # Имитируем команду расписание
        await message.answer(f"📅 Загружаю расписание для группы {user.group}...")
        
        # Используем логику из echo.py
        now = datetime.now()
        
        if now.month >= 9:
            academic_year = f"{now.year}/{now.year + 1}"
        else:
            academic_year = f"{now.year - 1}/{now.year}"
        
        half_year = 1 if now.month in [9, 10, 11, 12, 1] else 2
        
        group_schedules = await schedule_service.get_group_schedule(
            group_name=user.group,
            academic_year=academic_year,
            half_year=half_year
        )
        
        if group_schedules:
            response_text = f"📅 Расписание группы {user.group}\n"
            for gs in group_schedules:
                response_text += f"\n📚 Семестр: {gs.semester}\n"
                response_text += f"🔗 Изображение: {gs.schedule_img_url}\n"
                response_text += f"📅 Дата добавления: {gs.created_at}\n"
                response_text += "─" * 20
            
            if group_schedules[0].schedule_img_url:
                await message.answer_photo(
                    photo=group_schedules[0].schedule_img_url,
                    caption=response_text[:1024]
                )
            else:
                await message.answer(response_text)
        else:
            await message.answer(f"Расписание для группы {user.group} не найдено")
    else:
        # Если не студент или нет группы, показываем меню
        from app.keyboards.user_keyboards import get_schedule_menu_keyboard
        await message.answer(
            "📅 Расписание\n\nВыберите действие:",
            reply_markup=get_schedule_menu_keyboard()
        ) 