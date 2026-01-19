from aiogram import Router, types
from aiogram import F
from app.models import UserResponse
from app.keyboards.user_keyboards import get_main_menu_keyboard, get_schedule_menu_keyboard, get_ticket_menu_keyboard
import logging

router = Router()
logger = logging.getLogger(__name__)


@router.message(F.text == "📋 Профиль")
async def profile_button(message: types.Message, user: UserResponse):
    """Обработчик кнопки Профиль"""
    role_emoji = {
        'student': '👨‍🎓',
        'teacher': '👨‍🏫',
        'admin': '👑',
        'moderator': '🛡️'
    }

    profile_text = ""

    for role in user.roles_list:
        emoji = role_emoji.get(role, '👤')

        profile_text += (
            f"{emoji} *Ваш профиль*\n\n"
            f"👤 *Имя:* {user.full_name}\n"
            f"🆔 *ID:* {user.telegram_id}\n"
            f"📧 *Username:* @{user.username if user.username else 'нет'}\n"
            f"🎓 *Роль:* {role}\n"
        )

    if user.group:
        profile_text += f"📚 *Группа:* {user.group}\n"

    profile_text += f"📅 *Дата регистрации:* {user.created_at}"

    await message.answer(
        profile_text,
        parse_mode="Markdown",
        reply_markup=get_main_menu_keyboard()
    )


@router.message(F.text == "📅 Расписание")
async def schedule_button(message: types.Message, user: UserResponse):
    """Обработчик кнопки Расписание"""
    if 'student' in user.roles_list and user.group:
        # Для студентов показываем расписание их группы
        await message.answer(
            f"📅 *Расписание группы {user.group}*\n\n"
            "Выберите период:",
            parse_mode="Markdown",
            reply_markup=get_schedule_menu_keyboard()
        )
    elif 'teacher' in user.roles_list:
        # Для преподавателей можно сделать выбор группы
        await message.answer(
            "📅 *Расписание*\n\n"
            "Выберите группу или период:",
            parse_mode="Markdown",
            reply_markup=get_schedule_menu_keyboard()
        )
    else:
        await message.answer(
            "📅 *Расписание*\n\n"
            "Выберите период:",
            parse_mode="Markdown",
            reply_markup=get_schedule_menu_keyboard()
        )


@router.message(F.text == "🎫 Тикеты")
async def tickets_button(message: types.Message, user: UserResponse):
    """Обработчик кнопки Тикеты"""
    await message.answer(
        "🎫 *Система тикетов*\n\n"
        "Здесь вы можете создать тикет для решения проблем "
        "или задать вопросы администрации.",
        parse_mode="Markdown",
        reply_markup=get_ticket_menu_keyboard()
    )


@router.message(F.text == "⚙️ Настройки")
async def settings_button(message: types.Message, user: UserResponse):
    """Обработчик кнопки Настройки"""
    await message.answer(
        "⚙️ *Настройки*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_main_menu_keyboard()
    )
