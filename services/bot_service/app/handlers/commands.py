from aiogram import Router, types
from aiogram.filters import Command
from aiogram.types import Message
from app.models import UserResponse
from keyboards.user_keyboards import get_main_menu_keyboard
import logging

logger = logging.getLogger(__name__)

router = Router()


@router.message(Command("menu"))
async def cmd_menu(message: Message, user: UserResponse):
    """Обработчик команды /menu"""
    await message.answer(
        "🏠 *Главное меню*\n\n"
        "Выберите действие:",
        parse_mode="Markdown",
        reply_markup=get_main_menu_keyboard()
    )


@router.message(Command("profile"))
async def cmd_profile(message: Message, user: UserResponse):
    """Обработчик команды /profile"""
    role_emoji = {
        'student': '👨‍🎓',
        'teacher': '👨‍🏫',
        'admin': '👑',
        'moderator': '🛡️'
    }
    
    emoji = role_emoji.get(user.role, '👤')
    
    profile_text = (
        f"{emoji} *Ваш профиль*\n\n"
        f"👤 *Имя:* {user.full_name}\n"
        f"🆔 *ID:* {user.telegram_id}\n"
        f"📧 *Username:* @{user.username if user.username else 'нет'}\n"
        f"🎓 *Роль:* {user.role}\n"
    )
    
    if user.group_name:
        profile_text += f"📚 *Группа:* {user.group_name}\n"
    
    profile_text += f"📅 *Дата регистрации:* {user.created_at}"
    
    await message.answer(
        profile_text,
        parse_mode="Markdown",
        reply_markup=get_main_menu_keyboard()
    )


@router.message(Command("schedule"))
async def cmd_schedule(message: Message, user: UserResponse):
    """Обработчик команды /schedule"""
    from keyboards.user_keyboards import get_schedule_menu_keyboard
    
    if user.role == 'student' and user.group_name:
        await message.answer(
            f"📅 *Расписание группы {user.group_name}*\n\n"
            "Выберите период:",
            parse_mode="Markdown",
            reply_markup=get_schedule_menu_keyboard()
        )
    elif user.role == 'teacher':
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


@router.message(Command("ticket"))
async def cmd_ticket(message: Message, user: UserResponse):
    """Обработчик команды /ticket - система тикетов"""
    """
    TODO: Полная реализация системы тикетов с PostgreSQL
    
    Планируемая структура таблицы tickets:
    
    CREATE TABLE tickets (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        ticket_type VARCHAR(50) NOT NULL,
        status VARCHAR(20) DEFAULT 'open',
        priority VARCHAR(20) DEFAULT 'medium',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        resolved_at TIMESTAMP NULL
    );
    
    Пример кода для создания тикета:
    
    # from app.database import async_session
    # from app.models.ticket import Ticket
    # from sqlalchemy import insert
    # 
    # async with async_session() as session:
    #     stmt = insert(Ticket).values(
    #         user_id=user.id,
    #         title=title,
    #         description=description,
    #         ticket_type=ticket_type,
    #         status='open',
    #         priority=priority
    #     )
    #     await session.execute(stmt)
    #     await session.commit()
    # 
    # Пример кода для получения тикетов:
    # 
    # from sqlalchemy import select
    # 
    # async with async_session() as session:
    #     stmt = select(Ticket).where(Ticket.user_id == user.id).order_by(Ticket.created_at.desc())
    #     result = await session.execute(stmt)
    #     tickets = result.scalars().all()
    """
    
    from keyboards.user_keyboards import get_ticket_menu_keyboard
    
    await message.answer(
        "🎫 *Система тикетов*\n\n"
        "*В разработке:*\n"
        "• Создание тикетов с приоритетами\n"
        "• Прикрепление файлов\n"
        "• Общение с поддержкой\n"
        "• История тикетов\n\n"
        "Скоро здесь будет полноценная система поддержки!",
        parse_mode="Markdown",
        reply_markup=get_ticket_menu_keyboard()
    )