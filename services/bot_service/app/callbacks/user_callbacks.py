from aiogram import Router, types
from aiogram.filters import Text
from app.models import UserResponse
from app.services.schedule_service import schedule_service
from keyboards.user_keyboards import get_main_menu_keyboard, get_schedule_menu_keyboard, get_ticket_menu_keyboard, get_ticket_types_keyboard, get_back_to_menu_keyboard
import logging
from datetime import datetime

router = Router()
logger = logging.getLogger(__name__)


@router.callback_query(Text("back_to_menu"))
async def back_to_menu(callback: types.CallbackQuery, user: UserResponse):
    """Возврат в главное меню"""
    await callback.message.edit_text(
        "🏠 *Главное меню*\n\n"
        "Выберите действие:",
        parse_mode="Markdown"
    )
    await callback.message.edit_reply_markup(reply_markup=get_main_menu_keyboard())
    await callback.answer()


@router.callback_query(Text("back_to_tickets"))
async def back_to_tickets(callback: types.CallbackQuery, user: UserResponse):
    """Возврат в меню тикетов"""
    await callback.message.edit_text(
        "🎫 *Система тикетов*\n\n"
        "Здесь вы можете создать тикет для решения проблем "
        "или задать вопросы администрации.",
        parse_mode="Markdown",
        reply_markup=get_ticket_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_today"))
async def schedule_today(callback: types.CallbackQuery, user: UserResponse):
    """Расписание на сегодня"""
    try:
        # Получаем текущую дату для определения учебного года и семестра
        now = datetime.now()
        
        # Определяем учебный год (примерная логика)
        academic_year = f"{now.year}-{now.year + 1}"
        
        # Определяем семестр (1 - сентябрь-январь, 2 - февраль-июнь)
        half_year = 1 if 9 <= now.month <= 12 or now.month == 1 else 2
        
        if user.role == 'student' and user.group_name:
            schedule = await schedule_service.get_group_schedule(
                group_name=user.group_name,
                academic_year=academic_year,
                half_year=half_year
            )
            
            if schedule:
                # Форматируем расписание
                schedule_text = f"📅 *Расписание на сегодня ({now.strftime('%d.%m.%Y')})*\n\n"
                for day_schedule in schedule:
                    # Фильтруем занятия на сегодня
                    schedule_date = datetime.strptime(day_schedule.date, "%Y-%m-%d").date()
                    if schedule_date == now.date():
                        schedule_text += f"*{day_schedule.day_of_week}:*\n"
                        for lesson in day_schedule.lessons:
                            schedule_text += f"🕐 {lesson.time}: {lesson.subject}\n"
                            if lesson.teacher:
                                schedule_text += f"   👨‍🏫 {lesson.teacher}\n"
                            if lesson.classroom:
                                schedule_text += f"   🏫 {lesson.classroom}\n"
                        schedule_text += "\n"
                
                if schedule_text == f"📅 *Расписание на сегодня ({now.strftime('%d.%m.%Y')})*\n\n":
                    schedule_text += "На сегодня занятий нет 🎉"
            else:
                schedule_text = "Не удалось загрузить расписание. Попробуйте позже."
        else:
            schedule_text = "Информация о расписании доступна только студентам с указанной группой."
        
        await callback.message.edit_text(
            schedule_text,
            parse_mode="Markdown",
            reply_markup=get_back_to_menu_keyboard()
        )
        
    except Exception as e:
        logger.error(f"Error getting schedule: {e}")
        await callback.message.edit_text(
            "Произошла ошибка при получении расписания.",
            reply_markup=get_back_to_menu_keyboard()
        )
    
    await callback.answer()


@router.callback_query(Text("schedule_tomorrow"))
async def schedule_tomorrow(callback: types.CallbackQuery, user: UserResponse):
    """Расписание на завтра"""
    await callback.message.edit_text(
        "📅 *Расписание на завтра*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_week"))
async def schedule_week(callback: types.CallbackQuery, user: UserResponse):
    """Расписание на неделю"""
    await callback.message.edit_text(
        "📅 *Расписание на неделю*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("schedule_month"))
async def schedule_month(callback: types.CallbackQuery, user: UserResponse):
    """Расписание на месяц"""
    await callback.message.edit_text(
        "📅 *Расписание на месяц*\n\n"
        "Функция в разработке...",
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("create_ticket"))
async def create_ticket(callback: types.CallbackQuery, user: UserResponse):
    """Создание тикета"""
    await callback.message.edit_text(
        "🎫 *Создание тикета*\n\n"
        "Выберите тип проблемы:",
        parse_mode="Markdown",
        reply_markup=get_ticket_types_keyboard()
    )
    await callback.answer()


@router.callback_query(Text("my_tickets"))
async def my_tickets(callback: types.CallbackQuery, user: UserResponse):
    """Просмотр моих тикетов"""
    """
    TODO: Реализация просмотра тикетов из PostgreSQL
    
    Пример кода для работы с БД:
    
    # from app.database import async_session
    # from app.models.ticket import Ticket
    # from sqlalchemy import select
    
    # async with async_session() as session:
    #     stmt = select(Ticket).where(Ticket.user_id == user.id)
    #     result = await session.execute(stmt)
    #     tickets = result.scalars().all()
    #     
    #     if tickets:
    #         ticket_list = "🎫 *Ваши тикеты:*\n\n"
    #         for ticket in tickets:
    #             status_emoji = "🟢" if ticket.status == 'open' else "🟡" if ticket.status == 'in_progress' else "🔴"
    #             ticket_list += f"{status_emoji} *{ticket.title}*\n"
    #             ticket_list += f"📅 {ticket.created_at}\n"
    #             ticket_list += f"📝 {ticket.description[:50]}...\n\n"
    #     else:
    #         ticket_list = "У вас пока нет созданных тикетов."
    """
    
    ticket_list = "🎫 *Ваши тикеты*\n\n"
    ticket_list += "Функция в разработке...\n\n"
    ticket_list += "Скоро вы сможете просматривать свои тикеты здесь."
    
    await callback.message.edit_text(
        ticket_list,
        parse_mode="Markdown",
        reply_markup=get_back_to_menu_keyboard()
    )
    await callback.answer()


@router.callback_query(Text(startswith="ticket_type_"))
async def select_ticket_type(callback: types.CallbackQuery, user: UserResponse):
    """Выбор типа тикета"""
    ticket_type = callback.data.replace("ticket_type_", "")
    
    ticket_types = {
        'tech_issue': "🚀 Техническая проблема",
        'schedule_question': "📚 Вопрос по расписанию",
        'group_issue': "👥 Проблема с группой",
        'other': "❓ Другое"
    }
    
    await callback.message.edit_text(
        f"Вы выбрали: *{ticket_types.get(ticket_type, 'Неизвестный тип')}*\n\n"
        "Введите описание проблемы:",
        parse_mode="Markdown"
    )
    
    # TODO: Здесь можно запустить FSM для сбора данных тикета
    # from aiogram.fsm.state import StatesGroup, State
    # from aiogram.fsm.context import FSMContext
    # 
    # class TicketCreation(StatesGroup):
    #     waiting_for_description = State()
    #     waiting_for_priority = State()
    # 
    # await state.set_state(TicketCreation.waiting_for_description)
    # await state.update_data(ticket_type=ticket_type)
    
    await callback.answer() 