from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ReplyKeyboardMarkup, KeyboardButton
from aiogram.utils.keyboard import InlineKeyboardBuilder, ReplyKeyboardBuilder


def get_main_menu_keyboard() -> ReplyKeyboardMarkup:
    """Клавиатура главного меню"""
    builder = ReplyKeyboardBuilder()
    
    builder.add(KeyboardButton(text="📋 Профиль"))
    builder.add(KeyboardButton(text="📅 Расписание"))
    builder.add(KeyboardButton(text="⚙️ Настройки"))
    
    builder.adjust(2, 1)
    return builder.as_markup(resize_keyboard=True)


def get_schedule_menu_keyboard() -> InlineKeyboardMarkup:
    """Инлайн клавиатура для меню расписания"""
    builder = InlineKeyboardBuilder()
    
    builder.add(
        InlineKeyboardButton(
            text="📅 На сегодня",
            callback_data="schedule_today"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="📅 На завтра",
            callback_data="schedule_tomorrow"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="📅 На неделю",
            callback_data="schedule_week"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="⚙️ Настройки группы",
            callback_data="group_settings"
        )
    )
    
    builder.adjust(2, 1, 1)
    return builder.as_markup()


def get_group_selection_keyboard(groups: list) -> InlineKeyboardMarkup:
    """Клавиатура для выбора группы"""
    builder = InlineKeyboardBuilder()
    
    for group in groups:
        builder.add(
            InlineKeyboardButton(
                text=group,
                callback_data=f"select_group_{group}"
            )
        )
    
    builder.adjust(2)
    return builder.as_markup()


def get_back_to_menu_keyboard() -> InlineKeyboardMarkup:
    """Кнопка возврата в меню"""
    builder = InlineKeyboardBuilder()
    
    builder.add(
        InlineKeyboardButton(
            text="🔙 В меню",
            callback_data="back_to_menu"
        )
    )
    
    return builder.as_markup()