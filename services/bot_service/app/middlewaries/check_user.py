from aiogram import BaseMiddleware
from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup, Message
from aiogram.fsm.context import FSMContext

from app.services.user_service import user_service


def get_registration_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(
                    text="Зарегистрироваться", callback_data="register_user")
            ]
        ]
    )
    return keyboard


class CheckUserMiddleware(BaseMiddleware):
    async def __call__(self, handler, event: Message, data: dict):
        state: FSMContext = data.get('state')
        if state:
            current_state = await state.get_state()
            if current_state == 'RegisterState:waiting_for_code':
                return await handler(event, data)

        user_id = event.from_user.id if event.from_user else 1
        user = await user_service.get_user(user_id)

        if user:
            data['user'] = user
            return await handler(event, data)
        else:
            data['user'] = None
            if isinstance(event, Message):
                await event.answer(
                    text="Привет! Для использования бота вам необходимо зарегистрироваться:",
                    reply_markup=get_registration_keyboard()
                )
