from datetime import datetime
import logging
from aiogram import Router
from aiogram.types import InputMediaPhoto, Message
from aiogram import F

from app.models import UserResponse
from app.services.n8n_service import n8n_service
from app.services.schedule_service import schedule_service

echo_router = Router()


@echo_router.message(F.text.startswith('расписание'))
async def group_schedule(message: Message, user: UserResponse):
    try:
        req = message.text.split(" ")

        if len(req) < 2:
            group_name = user.group
        else:
            group_name = req[1]

        group_schedules = await schedule_service.get_group_schedule(
            group_name=group_name,
        )

        if group_schedules:
            media = [InputMediaPhoto(media=gs.schedule_img_url)
                     for gs in group_schedules]
            await message.answer_media_group(media=media)
        else:
            response_text = "ошибка"
            await message.answer(response_text)

    except Exception:
        response_text = """
Неверный формат запроса
Необходимо отправлять запрос в следующем виде:

    расписание группа
    расписание СА-501
        """

        await message.answer(response_text)


@echo_router.message(F.text.startswith('изменения'))
async def schedule_changes(message: Message, user: UserResponse):
    try:
        params = message.text.split(' ')
        if len(params) == 2:
            date = datetime.strptime(params[1], '%Y-%m-%d')
        else:
            date = None

        resp = await schedule_service.get_schedule_changes(date)

        if resp:
            media = [InputMediaPhoto(media=url) for url in resp.image_urls]
            await message.answer(f"""
Изменения на  {resp.date.strftime('%d.%m.%Y')}
Коментарий: {resp.description}
""")
            await message.answer_media_group(media=media)
            return

        await message.answer("Изменений в расписании не найдено")

    except Exception as e:
        logging.debug(f"Failed handle changes: {e}")
        await message.answer("Ошибка при получении изменений расписания")


@echo_router.message(F.text)
async def echo_handler(message: Message):
    msg = await message.answer("Подождите, идет обработка запроса")
    if message.text:
        ai_response = await n8n_service.proccess(message.text)
        if not ai_response:
            response_text = "ошибка обработки"
        else:
            response_text = ai_response.text
            media = [InputMediaPhoto(media=url)
                     for url in ai_response.photo_urls]
            await msg.edit_text(response_text)
            await message.answer_media_group(media=media)
            return
    else:
        response_text = "ошибка"
    await msg.edit_text(response_text)
