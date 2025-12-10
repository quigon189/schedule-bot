from aiogram import Router, types
from aiogram.filters import Text
from services.auth_service import AuthService
from services.schedule_service import ScheduleService
from keyboards.user_keyboards import get_main_menu_keyboard, get_schedule_menu_keyboard, get_back_to_menu_keyboard
import logging

router = Router()
auth_service = AuthService()
schedule_service = ScheduleService()

async def IsRegistred(telegram_id: int) -> bool:
    """Проверка регистрации пользователя"""
    user = await auth_service.get_user(telegram_id)
    return user is not None


@router.callback_query(Text("back_to_menu"))
async def back_to_menu(callback: types.CallbackQuery):
    """Возврат в главное меню"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await callback.message.edit_text(
        "Главное меню:",
        reply_markup=get_main_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_today"))
async def schedule_today(callback: types.CallbackQuery):
    """Расписание на сегодня"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    # Здесь будет логика получения расписания на сегодня
    await callback.message.edit_text(
        "📅 *Расписание на сегодня*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_tomorrow"))
async def schedule_tomorrow(callback: types.CallbackQuery):
    """Расписание на завтра"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await callback.message.edit_text(
        "📅 *Расписание на завтра*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_week"))
async def schedule_week(callback: types.CallbackQuery):
    """Расписание на неделю"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await callback.message.edit_text(
        "📅 *Расписание на неделю*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("group_settings"))
async def group_settings(callback: types.CallbackQuery):
    """Настройки группы"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    await callback.message.edit_text(
        "⚙️ *Настройки группы*\n\n"
        "Выберите вашу группу:",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text(startswith="select_group_"))
async def select_group(callback: types.CallbackQuery):
    """Выбор группы"""
    telegram_id = callback.from_user.id
    
    # Проверяем регистрацию
    if not await IsRegistred(telegram_id):
        await callback.message.answer("Вы не зарегистрированы. Используйте /start для регистрации.")
        return
    
    group_name = callback.data.replace("select_group_", "")
    
    await callback.message.edit_text(
        f"✅ Группа {group_name} выбрана!",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()