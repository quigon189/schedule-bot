from aiogram import F
from aiogram.types import Message
from aiogram.filters import Command
from aiogram import Router
from callbacks.user_callbacks import schedule_callback_changes , schedule_callback_main
from app.services.schedule_service import ScheduleService
from keyboards.user_keyboards import schedule_menu , profile_menu

router = Router(name="user_router")

@router.message(Command("schedule"))
async def schedule(message: Message):
    keyboard = schedule_menu()
    await message.answer("Что Вам необходимо?", reply_markup = keyboard)

@router.message(Command("profile"))
async def profile(message: Message):
    keyboard = profile_menu()
    await message.answer("Выберите действие:", reply_markup= keyboard)