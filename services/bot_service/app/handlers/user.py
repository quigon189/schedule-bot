from aiogram import Router, types
from aiogram import F
from app.models import UserResponse
from app.keyboards.user_keyboards import get_main_menu_keyboard, get_schedule_menu_keyboard, get_ticket_menu_keyboard
import logging
from datetime import datetime
from app.services.schedule_service import schedule_service

router = Router()
logger = logging.getLogger(__name__)

@router.message(F.text == "📋 Профиль")
async def profile_button(message: types.Message, user: UserResponse):
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

@router.message(F.text == "📅 Расписание")
async def schedule_button(message: types.Message, user: UserResponse):

    await message.answer("Сначала зарегистрируйтесь!")
    return
    
    await message.answer(
        "📅 Расписание\n\nВыберите действие:",
        reply_markup=get_schedule_menu_keyboard()
    )

@router.message(F.text == "🎫 Тикеты")
async def tickets_button(message: types.Message, user: UserResponse):

    await message.answer("Сначала зарегистрируйтесь!")
    return
    
    await message.answer(
        "🎫 Система тикетов\n\nФункция в разработке...",
        reply_markup=get_ticket_menu_keyboard()
    )

@router.message(F.text == "⚙️ Настройки")
async def settings_button(message: types.Message, user: UserResponse):

    await message.answer("Сначала зарегистрируйтесь!")
    return
    
    await message.answer(
        "⚙️ Настройки\n\nФункция в разработке...",
        reply_markup=get_main_menu_keyboard()
    )

@router.message(F.text == "📅 Расписание")
async def schedule_button(message: types.Message, user: UserResponse):
    await message.answer("Сначала зарегистрируйтесь!")
    return
    
    # Аналогичная логика как в команде /schedule
    if hasattr(user, 'roles_list') and 'student' in user.roles_list and user.group:
        await message.answer(f"📅 Загружаю расписание для группы {user.group}...")
        
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
        await message.answer(
            "📅 Расписание\n\nВыберите действие:",
            reply_markup=get_schedule_menu_keyboard()
        )