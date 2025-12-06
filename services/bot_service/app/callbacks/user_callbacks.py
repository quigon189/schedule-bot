from aiogram.types import CallbackQuery
from aiogram import F
from aiogram.utils.callback_answer import CallbackAnswer
from aiogram.filters.callback_data import CallbackData
from app.services.schedule_service import ScheduleService
from handlers.user import router
from keyboards.common_keyboards import start_menu

@router.callback_query(F.data == "schedule_changes")
async def schedule_callback_changes(callback: CallbackQuery):
    await ScheduleService.get_group_schedule()
    await callback.answer()

@router.callback_query(F.data == "schedule_main")
async def schedule_callback_main(callback: CallbackQuery):
    await ScheduleService.get_group_schedule()
    await callback.answer()

@router.callback_query(F.data == "change_profile")
async def ch_prof(callback: CallbackQuery):
    await callback.answer()

@router.callback_query(F.data == "back_to_main")
async def back_to_main(callback: CallbackQuery):
    keyboard = start_menu()
    await callback.message.edit_text(
        "Главное меню:",
        reply_markup=keyboard
    )
    await callback.answer()