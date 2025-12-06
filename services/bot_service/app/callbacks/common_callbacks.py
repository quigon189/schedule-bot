from aiogram.types import CallbackQuery
from aiogram import F
from aiogram.utils.callback_answer import CallbackAnswer
from aiogram.filters.callback_data import CallbackData
from keyboards.user_keyboards import schedule_menu, profile_menu
from handlers.common import router
from keyboards.common_keyboards import start_menu

@router.callback_query(F.data == "profile_menu")
async def profile_menu(callback: CallbackQuery):
    keyboard = profile_menu
    await callback.message.edit_text(
        "Меню профилей:",
        reply_markup=keyboard
    )
    await callback.anwser()

@router.callback_query(F.data == "schedule_menu")
async def schedule_menu(callback: CallbackQuery):
    keyboard = schedule_menu
    await callback.message.edit_text(
        "Меню расписаний:",
        reply_markup=keyboard
    )
    await callback.anwser()

@router.callback_query(F.data == "back_to_main")
async def back_to_main(callback: CallbackQuery):
    keyboard = start_menu()
    await callback.message.edit_text(
        "Главное меню:",
        reply_markup=keyboard
    )
    await callback.answer()