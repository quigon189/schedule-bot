from aiogram import Router
from aiogram.types import Message
from aiogram import F

from app.services.schedule_service import schedule_service
from app.services.user_service import user_service

echo_router = Router()


@echo_router.message(F.text.startswith('расписание'))
async def group_schedule(message: Message):
    try:
        req = message.text.split(" ")

        group_schedules = await schedule_service.get_group_schedule(
            group_name=req[1],
            academic_year=req[2],
            half_year=int(req[3])
        )

        if group_schedules:
            response_text = ""
            for gs in group_schedules:
                response_text += f"""
Расписание группы {gs.group_name}:
Cеместр: {gs.semester}

{gs.schedule_img_url}

Дата добавлния {gs.created_at}
                """
        else:
            response_text = "ошибка"

        await message.answer(response_text)

    except Exception:
        response_text = """
Неверный формат запроса

Необходимо отправлять запрос в следующем виде:

    расписание группа учебный-год полугодие
    расписание СА-501 2025/2026 1
        """

        await message.answer(response_text)


@echo_router.message(F.text)
async def echo_handler(message: Message):
    user = await user_service.get_user(
        user_id=message.from_user.id
    )

    if user:
        response_text = f"""
📨 Вы написали: {message.text}

Вы зарегистрированны в системе
👤 Ваши данные:
Username: {user.username}
FullName: {user.full_name}
        """
    else:
        response_text = f"""
📨 Вы написали: {message.text}

Вы НЕ зарегистрированны в системе
        """

    await message.answer(response_text)
