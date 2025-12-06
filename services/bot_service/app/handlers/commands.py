from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

from app.models import TelegramUser
from app.services.user_service import user_service

from app.keyboards.common_keyboards import start_menu

command_router = Router()


@command_router.message(Command("help"))
async def cmd_help(message: Message):
    help_text = """
🤖 Доступные команды:

/start - Начать работу с ботом
/help - Показать эту справку

Бот отвечает на ваши сообщения эхом.
    """

    await message.answer(help_text)


@command_router.message(Command("start"))
async def cmd_start(message: Message):
    user = TelegramUser(
        id=message.from_user.id,
        username=message.from_user.username,
        first_name=message.from_user.first_name,
        last_name=message.from_user.last_name
    )

    registered_user = await user_service.register_user(user)

    if registered_user:
        welcome_text = f"""
✅ Добро пожаловать, {user.full_name}!

Вы успешно зарегистрированы в системе.

Теперь вы можете использовать все возможности бота!
        """
        keyboard = start_menu(),
        await message.answer("Добро пожаловать! Доступные функции:", reply_markup= keyboard)
    else:
        welcome_text = f"""
👋 Привет, {user.full_name}!

К сожалению, не удалось завершить регистрацию.
Попробуйте позже или обратитесь к администратору.
        """

    await message.answer(welcome_text)
