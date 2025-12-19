from aiogram import BaseMiddleware
from aiogram.types import InlineKeyboardButton, InlineKeyboardMarkup, Message, CallbackQuery
from aiogram.fsm.context import FSMContext

from app.services.user_service import user_service

def get_registration_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(
                    text="Зарегистрироваться", 
                    callback_data="register_user"
                )
            ]
        ]
    )
    return keyboard

class CheckUserMiddleware(BaseMiddleware):
    async def __call__(self, handler, event: Message | CallbackQuery, data: dict):
        # Получаем user_id из события
        if isinstance(event, Message):
            user_id = event.from_user.id
        elif isinstance(event, CallbackQuery):
            user_id = event.from_user.id
        else:
            return await handler(event, data)

        # Проверяем, есть ли состояние регистрации
        state: FSMContext = data.get('state')
        if state:
            current_state = await state.get_state()
            # Пропускаем если пользователь в процессе регистрации
            if current_state and 'register' in current_state.lower():
                return await handler(event, data)

        # Получаем пользователя из auth_service через user_service
        user = await user_service.get_user(user_id)
        
        if user and user.is_active:
            data['user'] = user
            return await handler(event, data)
        else:
            # Если пользователь не найден или не активен
            data['user'] = None
            
            if isinstance(event, Message):
                # Не показываем регистрацию на команды /start и /help
                if event.text and (event.text.startswith('/start') or event.text.startswith('/help')):
                    return await handler(event, data)
                
                await event.answer(
                    text="🔒 Для использования бота необходимо зарегистрироваться.\n\n"
                         "Получите код регистрации у администратора.",
                    reply_markup=get_registration_keyboard()
                )
                return
            elif isinstance(event, CallbackQuery):
                # Если коллбэк не для регистрации
                if event.data != "register_user":
                    await event.answer(
                        "Сначала зарегистрируйтесь!",
                        show_alert=True
                    )
                    return
                else:
                    return await handler(event, data)