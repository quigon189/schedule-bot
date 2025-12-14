from aiogram.types import InlineKeyboardMarkup, InlineKeyboardButton, ReplyKeyboardMarkup, KeyboardButton
from aiogram.utils.keyboard import InlineKeyboardBuilder, ReplyKeyboardBuilder


def get_main_menu_keyboard() -> ReplyKeyboardMarkup:
    """Клавиатура главного меню"""
    builder = ReplyKeyboardBuilder()

    builder.add(KeyboardButton(text="📋 Профиль"))
    builder.add(KeyboardButton(text="📅 Расписание"))
    builder.add(KeyboardButton(text="🎫 Тикеты"))
    builder.add(KeyboardButton(text="⚙️ Настройки"))

    builder.adjust(2)
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
            text="📅 На месяц",
            callback_data="schedule_month"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="🔙 Назад",
            callback_data="back_to_menu"
        )
    )

    builder.adjust(2, 2, 1)
    return builder.as_markup()


def get_ticket_menu_keyboard() -> InlineKeyboardMarkup:
    """Инлайн клавиатура для меню тикетов"""
    builder = InlineKeyboardBuilder()

    builder.add(
        InlineKeyboardButton(
            text="🎫 Создать тикет",
            callback_data="create_ticket"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="📋 Мои тикеты",
            callback_data="my_tickets"
        )
    )
    builder.add(
        InlineKeyboardButton(
            text="🔙 Назад",
            callback_data="back_to_menu"
        )
    )

    builder.adjust(1, 1, 1)
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


def get_ticket_types_keyboard() -> InlineKeyboardMarkup:
    """Клавиатура для выбора типа тикета"""
    builder = InlineKeyboardBuilder()

    ticket_types = [
        ("🚀 Техническая проблема", "tech_issue"),
        ("📚 Вопрос по расписанию", "schedule_question"),
        ("👥 Проблема с группой", "group_issue"),
        ("❓ Другое", "other")
    ]

    for text, callback_data in ticket_types:
        builder.add(
            InlineKeyboardButton(
                text=text,
                callback_data=f"ticket_type_{callback_data}"
            )
        )

    builder.add(
        InlineKeyboardButton(
            text="🔙 Назад",
            callback_data="back_to_tickets"
        )
    )

    builder.adjust(1)
    return builder.as_markup()
