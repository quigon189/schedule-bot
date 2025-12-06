from aiogram import F
from aiogram.types import Message
from aiogram.filters import Command
from aiogram import Router
from keyboards.common_keyboards import start_menu
from app.services.auth_service import CheckIfRegistred

router = Router(name="common_router")

@router.message(Command("start"))
async def start_cmd(message: Message):
    if CheckIfRegistred.IsRegistred == True:
            keyboard = start_menu(),
            await message.answer("Добро пожаловать! Доступные функции:", reply_markup= keyboard)
    else:
       await message.answer("Для начала работы зарегистрируйтесь в системе!")

@router.message(Command("help"))
async def help_cmd(message: Message):
    keyboard = start_menu(),
    await message.anwser("В этом боте Вы можете просмотреть расписание для своей группы или действующие изменения в расписании. Другие функции пока в разработки. Для начала работы выберете один из пунктов ниже:", reply_markup= keyboard)

