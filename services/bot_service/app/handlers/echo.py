from datetime import datetime
import logging
from aiogram import Router
from aiogram.types import InputMediaPhoto, Message, CallbackQuery
from aiogram import F

from app.models import UserResponse
from app.services.schedule_service import schedule_service
from app.services.user_service import user_service
from app.keyboards.user_keyboards import get_main_menu_keyboard

echo_router = Router()

@echo_router.message(F.text.startswith('расписание'))
async def group_schedule(message: Message, user: UserResponse = None):
    if not user:
        await message.answer("Сначала зарегистрируйтесь!")
        return
    
    try:
        # Если пользователь студент и у него есть группа, показываем его расписание
        if hasattr(user, 'roles_list') and 'student' in user.roles_list and user.group:
            await message.answer(f"📅 Загружаю расписание для группы {user.group}...")
            
            # Определяем текущий учебный год и семестр
            now = datetime.now()
            
            # Учебный год (пример: если месяц >= 9, то учебный год текущий/следующий)
            if now.month >= 9:
                academic_year = f"{now.year}/{now.year + 1}"
            else:
                academic_year = f"{now.year - 1}/{now.year}"
            
            # Полугодие (1 - сентябрь-январь, 2 - февраль-июнь)
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
                
                # Если есть URL изображения, показываем его
                if group_schedules[0].schedule_img_url:
                    await message.answer_photo(
                        photo=group_schedules[0].schedule_img_url,
                        caption=response_text[:1024]  # Ограничение Telegram
                    )
                else:
                    await message.answer(response_text)
            else:
                await message.answer(f"Расписание для группы {user.group} не найдено")
                
        else:
            # Если не студент или нет группы, просим ввести запрос вручную
            await message.answer(
                "Введите запрос в формате:\n"
                "расписание [группа] [учебный год] [полугодие]\n"
                "Пример: расписание СА-501 2024/2025 1"
            )
            
    except Exception as e:
        logging.debug(f"Error in group_schedule: {e}")
        await message.answer(
            "Неверный формат запроса. Используйте:\n"
            "расписание [группа] [учебный год] [полугодие]\n"
            "Или просто 'расписание' для своей группы"
        )

@echo_router.message(F.text.startswith('изменения'))
async def schedule_changes(message: Message, user: UserResponse = None):
    if not user:
        await message.answer("Сначала зарегистрируйтесь!")
        return
    
    try:
        params = message.text.split(' ')
        date = None
        
        if len(params) >= 2:
            try:
                date = datetime.strptime(params[1], '%Y-%m-%d')
            except ValueError:
                date = None
        
        resp = await schedule_service.get_schedule_changes(date)
        
        if resp:
            # Если есть URL изображений
            if resp.image_urls:
                media = [InputMediaPhoto(media=url) for url in resp.image_urls]
                await message.answer(
                    f"📊 Изменения на {resp.date.strftime('%d.%m.%Y')}\n"
                    f"📝 {resp.description}"
                )
                await message.answer_media_group(media=media[:10])  # Ограничение: 10 фото
            else:
                await message.answer(
                    f"📊 Изменения на {resp.date.strftime('%d.%m.%Y')}\n"
                    f"📝 {resp.description}\n"
                    f"🖼️ Нет изображений"
                )
        else:
            date_text = f"на {date.strftime('%d.%m.%Y')}" if date else ""
            await message.answer(f"Изменений в расписании {date_text} не найдено")
            
    except Exception as e:
        logging.debug(f"Failed handle changes: {e}")
        await message.answer("Ошибка при получении изменений расписания")

@echo_router.callback_query(F.data == "schedule_today")
async def schedule_today_callback(callback: CallbackQuery, user: UserResponse = None):
    if not user:
        await callback.answer("Сначала зарегистрируйтесь!", show_alert=True)
        return
    
    # Просто показываем инструкцию
    await callback.message.edit_text(
        "📅 Расписание на сегодня\n\n"
        "Используйте команду /schedule\n"
        "Или напишите 'расписание' для своей группы\n\n"
        "Если вы студент и у вас указана группа, "
        "бот автоматически покажет ваше расписание.",
        reply_markup=get_main_menu_keyboard()
    )
    await callback.answer()

@echo_router.callback_query(F.data == "schedule_changes")
async def schedule_changes_callback(callback: CallbackQuery, user: UserResponse = None):
    if not user:
        await callback.answer("Сначала зарегистрируйтесь!", show_alert=True)
        return
    
    await callback.message.edit_text(
        "📊 Изменения расписания\n\n"
        "Напишите 'изменения' для просмотра всех изменений\n"
        "Или 'изменения 2024-12-25' для конкретной даты",
        reply_markup=get_main_menu_keyboard()
    )
    await callback.answer()

# Общие коллбэки остаются
@echo_router.callback_query(F.data.in_([
    "schedule_tomorrow", "schedule_week", "schedule_month",
    "create_ticket", "my_tickets"
]))
async def handle_other_callbacks(callback: CallbackQuery):
    await callback.message.edit_text("📝 Функция в разработке...")
    await callback.answer()

@echo_router.callback_query(F.data == "back_to_menu")
async def back_to_menu_callback(callback: CallbackQuery, user: UserResponse = None):
    if user:
        await callback.message.edit_text(
            "🏠 Главное меню\n\nВыберите действие:",
            reply_markup=get_main_menu_keyboard()
        )
    else:
        await callback.message.edit_text("Сначала зарегистрируйтесь!")
    await callback.answer()

@echo_router.message(F.text)
async def echo_handler(message: Message):
    user = await user_service.get_user(user_id=message.from_user.id)

    if user:
        response_text = f"""
📨 Вы написали: {message.text}

✅ Вы зарегистрированы в системе
👤 Ваши данные:
Имя: {user.full_name}
Группа: {user.group or 'не указана'}
        """
    else:
        response_text = f"""
📨 Вы написали: {message.text}

❌ Вы НЕ зарегистрированы в системе
Для регистрации получите код у администратора
        """

    await message.answer(response_text)