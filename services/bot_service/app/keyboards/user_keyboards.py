from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ReplyKeyboardMarkup, KeyboardButton
from aiogram.utils.keyboard import InlineKeyboardBuilder, ReplyKeyboardBuilder

def get_main_menu_keyboard() -> ReplyKeyboardMarkup:
    builder = ReplyKeyboardBuilder()
    builder.add(KeyboardButton(text="📋 Профиль"))
    builder.add(KeyboardButton(text="📅 Расписание"))
    builder.adjust(2)
    return builder.as_markup(resize_keyboard=True)

def get_schedule_menu_keyboard() -> InlineKeyboardMarkup:
    builder = InlineKeyboardBuilder()
    builder.add(InlineKeyboardButton(text="📅 Текущее", callback_data="schedule_current"))
    builder.add(InlineKeyboardButton(text="📊 Изменения на завтра", callback_data="schedule_changes"))
    builder.add(InlineKeyboardButton(text="🔙 Назад", callback_data="back_to_menu"))
    builder.adjust(2, 1)
    return builder.as_markup()

def get_back_to_menu_keyboard() -> InlineKeyboardMarkup:
    builder = InlineKeyboardBuilder()
    builder.add(InlineKeyboardButton(text="🔙 В меню", callback_data="back_to_menu"))
    return builder.as_markup()
    