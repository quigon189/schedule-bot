from aiogram import Router
from aiogram.filters import Command
from aiogram.types import Message
from app.models import UserResponse
from app.keyboards.user_keyboards import get_main_menu_keyboard
from app.services.schedule_service import schedule_service
from datetime import datetime
import logging

logger = logging.getLogger(__name__)

router = Router()

@router.message(Command("start"))
async def cmd_start(message: Message, user: UserResponse = None):
    if user:
        # Пользователь уже зарегистрирован
        await message.answer(
            f"✅ Добро пожаловать, {user.full_name}!\n\n"
            "Вы уже зарегистрированы в системе.\n"
            "Используйте меню или команды:\n"
            "/menu - главное меню\n"
            "/profile - ваш профиль\n"
            "/schedule - расписание\n"
            "/help - справка",
            reply_markup=get_main_menu_keyboard()
        )
    else:
        # Пользователь не зарегистрирован
        await message.answer(
            "👋 Добро пожаловать в систему расписаний!\n\n"
            "🔐 Для использования бота необходимо зарегистрироваться.\n\n"
            "📝 Получите код регистрации у администратора вашего учебного заведения "
            "или преподавателя, затем нажмите кнопку 'Зарегистрироваться'.\n\n"
            "❓ Если у вас возникли проблемы с регистрацией, "
            "обратитесь к администратору."
        )

@router.message(Command("menu"))
async def cmd_menu(message: Message, user: UserResponse):
    # Middleware гарантирует что user есть
    await message.answer(
        "🏠 Главное меню\n\nВыберите действие:",
        reply_markup=get_main_menu_keyboard()
    )

@router.message(Command("profile"))
async def cmd_profile(message: Message, user: UserResponse):
    # Формируем текст профиля
    roles_text = ", ".join(user.roles_list) if user.roles_list else "нет ролей"
    
    profile_text = (
        f"👤 Ваш профиль\n\n"
        f"📝 Имя: {user.full_name}\n"
        f"🆔 ID: {user.telegram_id}\n"
        f"📧 Username: @{user.username if user.username else 'нет'}\n"
        f"🎭 Роли: {roles_text}\n"
        f"📊 Статус: {'✅ Активен' if user.is_active else '❌ Неактивен'}\n"
    )
    
    if user.group:
        profile_text += f"📚 Группа: {user.group}\n"
    
    profile_text += f"📅 Регистрация: {user.created_at.strftime('%d.%m.%Y %H:%M')}\n"
    profile_text += f"🔄 Обновлен: {user.updated_at.strftime('%d.%m.%Y %H:%M')}"
    
    await message.answer(profile_text, reply_markup=get_main_menu_keyboard())

@router.message(Command("schedule"))
async def cmd_schedule(message: Message, user: UserResponse):
    # Проверяем, есть ли у пользователя группа
    if user.group:
        await message.answer(f"📅 Загружаю расписание для группы {user.group}...")
        
        # Определяем учебный год и семестр
        now = datetime.now()
        
        if now.month >= 9:
            academic_year = f"{now.year}/{now.year + 1}"
        else:
            academic_year = f"{now.year - 1}/{now.year}"
        
        half_year = 1 if now.month in [9, 10, 11, 12, 1] else 2
        
        # Получаем расписание через schedule_service
        group_schedules = await schedule_service.get_group_schedule(
            group_name=user.group,
            academic_year=academic_year,
            half_year=half_year
        )
        
        if group_schedules:
            response_text = f"📅 Расписание группы {user.group}\n"
            for gs in group_schedules:
                response_text += f"\n📚 Семестр: {gs.semester}\n"
                if gs.schedule_img_url:
                    response_text += f"🖼️ Изображение расписания\n"
                response_text += f"📅 Дата: {gs.created_at.strftime('%d.%m.%Y')}\n"
                response_text += "─" * 20
            
            # Если есть URL изображения, отправляем фото
            if group_schedules[0].schedule_img_url:
                try:
                    await message.answer_photo(
                        photo=group_schedules[0].schedule_img_url,
                        caption=response_text[:1024]
                    )
                except:
                    await message.answer(response_text)
            else:
                await message.answer(response_text)
        else:
            await message.answer(f"❌ Расписание для группы {user.group} не найдено")
    else:
        # У пользователя нет группы
        await message.answer(
            "📅 Расписание\n\n"
            "У вас не указана учебная группа в профиле.\n"
            "Обратитесь к администратору для добавления группы."
        )

@router.message(Command("help"))
async def cmd_help(message: Message, user: UserResponse = None):
    help_text = (
        "📚 Справка по командам:\n\n"
        "/start - Начало работы с ботом\n"
        "/menu - Главное меню\n"
        "/profile - Показать ваш профиль\n"
        "/schedule - Показать расписание\n"
        "/help - Эта справка\n\n"
    )
    
    if user and "admin" in user.roles_list:
        help_text += "⚙️ Команды администратора:\n"
        help_text += "/code - Создать код регистрации\n\n"
    
    help_text += (
        "📱 Также вы можете использовать кнопки меню:\n"
        "• 📋 Профиль - информация о вас\n"
        "• 📅 Расписание - ваше расписание\n"
        "• 🎫 Тикеты - система поддержки\n"
        "• ⚙️ Настройки - настройки бота\n\n"
        "❓ По вопросам обращайтесь к администратору."
    )
    
    await message.answer(help_text)

@router.message(Command("logout"))
async def cmd_logout(message: Message, user: UserResponse):
    # Очищаем кэш для этого пользователя
    from app.services.auth_service import auth_service
    await auth_service.invalidate_cache(message.from_user.id)
    
    await message.answer(
        "🔓 Кэш вашего профиля очищен.\n"
        "При следующем обращении данные будут загружены заново."
    )