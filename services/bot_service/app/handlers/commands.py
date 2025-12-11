from aiogram import Router
from aiogram.types import Message
from aiogram.filters import Command

from app.models import TelegramUser, UserResponse
from app.services.auth_service import auth_service
from app.services.user_service import user_service

command_router = Router()


@command_router.message(Command("help"))
async def cmd_help(message: Message):
    help_text = """
🤖 Доступные команды:

/start - Начать работу с ботом
/help - Показать эту справку
/code - Создать код регистрации

Бот отвечает на ваши сообщения эхом.
    """

    await message.answer(help_text)


@command_router.message(Command("start"))
async def cmd_start(message: Message):
    help_text = """
🤖 Доступные команды:

/start - Начать работу с ботом
/help - Показать эту справку

для админов
/code - Создать код регистрации

Бот отвечает на ваши сообщения эхом.
    """

    await message.answer(help_text)
