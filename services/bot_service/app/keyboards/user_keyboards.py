from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup
from aiogram.utils.keyboard import InlineKeyboardBuilder

def schedule_menu():
    builder = InlineKeyboardBuilder()
    builder.add(
        InlineKeyboardButton(text= "Основное расписание", callback_data="schedule_main"),
        InlineKeyboardButton(text= "Изменения в расписании", callback_data="schedule_changes"),
        InlineKeyboardButton(text="Главное меню", callback_data="back_to_main")
    )
    builder.adjust(2)  # 2 кнопки в строке
    return builder.as_markup()

def profile_menu():
    builder = InlineKeyboardBuilder()
    builder.add(
        InlineKeyboardButton(text="Редактировать", callback_data="change_profile"),
        InlineKeyboardButton(text="Главное меню", callback_data="back_to_main")
    )
    builder.adjust(2)  # 2 кнопки в строке
    return builder.as_markup()
