from datetime import datetime, timedelta
import re

from aiogram import Router, CallbackQuery
from aiogram.types import InputMediaPhoto
from aiogram.fsm.state import State, StatesGroup, FSMContext

from app.handlers import echo_router, admin_router, com_router, user_router, register_router
from app.servces.schedule_service import schedule_service


# REGISTER CALLBACKS

@register_router.callback_query(F.data == "register_user")
async def register_user(callback_query: CallbackQuery, state: FSMContext):
    await state.set_state(RegisterState.waiting_for_code)
    await callback_query.message.answer(
        "🔐 Введите код регистрации:\n\n"
        "Код должен быть получен у администратора или преподавателя.\n"
        "Введите код без пробелов и дополнительных символов."
    )
    await callback_query.answer()

# USER CALLBACKS

class RegisterState(StatesGroup):
    waiting_for_group = State()

def get_half_year(date: datetime) -> int:
    return 1 if 9 <= date.month <= 12 else 2

@echo_router.callback_query(F.data == "schedule_current")
async def schedule_today_callback(callback: CallbackQuery, user: UserResponse, state: FSMContext):
    # TODO: добавить определение учебного года, полугодия
    if "student" in user.roles_list:
        today = datetime.now()
        year_range = f"{today.year}/{today.year + 1}"
        group = user.group
        service_response = await schedule_service.get_group_schedule(
        group_name = user.group ,
        academic_year = year_range,
        half_year = get_half_year(today)
    )    
        if service_response:    
            await callback.message.answer_photo(service_response.schedule_img_url)
            await callback.message.answer(f"Текущее расписание для группы {user.group}")
        else:
            await callback.message.answer("Не удалось найти расписание.")
    else:
        await callback.message.answer("Введите название целевой группы в формате БУКВЫ-ЦИФРЫ. Пример:СА-501")
        await state.set_state(RegisterState.waiting_for_group)
    await callback.answer()

@echo_router.message(RegisterState.waiting_for_group)
async def process_group_input(message: Message, state: FSMContext, user: UserResponse):
    group = message.text.strip().upper()
    if not re.match(r"^[А-Я]{1,3}-\d{3}$", group):
        await message.answer(
            "Некорректный формат! Попробуйте ещё раз. \n"
            "Пример правильного формата: **СА-501**"
        )
        return
    today = datetime.now()
    year_range = f"{today.year}/{today.year + 1}"
    group = user.group
    service_response = await schedule_service.get_group_schedule(
    group_name = user.group ,
    academic_year = year_range,
    half_year = get_half_year(today)
    )    
    if service_response:    
        await callback.message.answer_photo(service_response.schedule_img_url)
        await callback.message.answer(f"Текущее расписание для группы {user.group}")
    else:
        await callback.message.answer("Не удалось найти расписание.")
    await state.clear()

@echo_router.callback_query(F.data == "schedule_changes")
async def schedule_changes_callback(callback: CallbackQuery, user: UserResponse):
    tommorrow = datetime.now() + timedelta(days = 1)
    service_response = await schedule_service.get_schedule_changes(
    date = tommorow.strftime("%Y-%m-%d")
    )    
    if service_response:    
        media = [InputMediaPhoto(media=url) for url in service_response.image_urls]
        await callback.message.anwser_media_group(media=media)
        await callback.message.answer("Изменения на завтра")
    else:
        await callback.message.answer("Не удалось найти изменения.")
    await callback.answer()

@echo_router.callback_query(F.data == "back_to_menu")
async def back_to_menu_callback(callback: CallbackQuery, user: UserResponse):
    await callback.message.edit_text(
        "🏠 Главное меню\n\nВыберите действие:",
        reply_markup=get_main_menu_keyboard()
    )
    await callback.answer()

# COMMANDS CALLBACKS

# ADMIN CALLBACKS

@admin_router.callback_query(F.data.startswith('group:'))
async def process_group_selection(callback_query: CallbackQuery):
    group = callback_query.data.split(':')[1]
    role = 'student'
    await create_final_code(callback_query, role, group)

async def create_final_code(callback_query: CallbackQuery, role: str, group: str = None):
    code = await auth_service.create_registration_code(
        role=role,
        group_name=group,
        max_uses=10,
        created_by=callback_query.from_user.id,
        expires=7 * 24 * 60 * 60  # 7 дней
    )
    
    if code:
        result_text = (
            f"✅ Код для {role} создан!\n\n"
            f"🔑 Код: `{code.code}`\n"
            f"👥 Группа: {group or 'любая'}\n"
            f"🔢 Использований: {code.max_uses}\n"
            f"⏰ Действует до: {code.expires_at.strftime('%d.%m.%Y %H:%M')}"
        )
    else:
        result_text = "❌ Ошибка при создании кода"
    
    await callback_query.message.edit_text(result_text)
    await callback_query.answer()

@admin_router.callback_query(F.data.startswith('create_code:'))
async def process_code_role(callback_query: CallbackQuery):
    role = callback_query.data.split(':')[1]
    
    if role == 'student':
        await callback_query.message.edit_text(
            'Выберите группу:',
            reply_markup=get_groups_keyboard()
        )
    else:
        await create_final_code(callback_query, role, None)
    
    await callback_query.answer()

@admin_router.message(Command("code"))
async def create_code(message: Message):
    await message.answer("Выберите роль для кода:", reply_markup=get_code_type_keyboard())

