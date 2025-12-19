from aiogram.types import CallbackQuery

@echo_router.callback_query(F.data == "schedule_today")
async def schedule_today_callback(callback: CallbackQuery, user: UserResponse):
    # Просто перенаправляем к основной функциональности
    await callback.message.answer(
        "📅 Используйте кнопку '📅 Расписание' в меню\n"
        "или команду /schedule для просмотра расписания"
    )
    await callback.answer()

@echo_router.callback_query(F.data == "schedule_changes")
async def schedule_changes_callback(callback: CallbackQuery, user: UserResponse):
    await callback.message.answer(
        "📊 Изменения расписания\n\n"
        "Напишите 'изменения' для просмотра изменений\n"
        "Или 'изменения 2024-12-25' для конкретной даты"
    )
    await callback.answer()

# Общие коллбэки
@echo_router.callback_query(F.data.in_([
    "schedule_tomorrow", "schedule_week", "schedule_month",
    "create_ticket", "my_tickets"
]))
async def handle_other_callbacks(callback: CallbackQuery, user: UserResponse):
    await callback.message.edit_text("📝 Функция в разработке...")
    await callback.answer()

@echo_router.callback_query(F.data == "back_to_menu")
async def back_to_menu_callback(callback: CallbackQuery, user: UserResponse):
    await callback.message.edit_text(
        "🏠 Главное меню\n\nВыберите действие:",
        reply_markup=get_main_menu_keyboard()
    )
    await callback.answer()