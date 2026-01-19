from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup
from aiogram.utils.keyboard import InlineKeyboardBuilder

def start_menu():
    builder = InlineKeyboardBuilder()
    builder.add(
        InlineKeyboardButton(text="Меню профиля", callback_data="schedule_menu"),
        InlineKeyboardButton(text="Меню расписаний", callback_data="profile_menu")
    )
    builder.adjust(2)  # 2 кнопки в строке
    return builder.as_markup()