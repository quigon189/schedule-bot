from aiogram import F, Router
from aiogram.types import CallbackQuery, Message
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.context import FSMContext
from app.models import TelegramUser
from app.services.user_service import user_service

register_router = Router()

class RegisterState(StatesGroup):
    waiting_for_code = State()

@register_router.callback_query(F.data == "register_user")
async def register_user(callback_query: CallbackQuery, state: FSMContext):
    await state.set_state(RegisterState.waiting_for_code)
    await callback_query.message.answer("Пожалуйста, введите код регистрации:")
    await callback_query.answer()

@register_router.message(RegisterState.waiting_for_code)
async def process_register_code(message: Message, state: FSMContext):
    code = message.text.strip()
    
    user = TelegramUser(
        id=message.from_user.id,
        username=message.from_user.username,
        first_name=message.from_user.first_name,
        last_name=message.from_user.last_name
    )

    reg_user = await user_service.register_user(code, user)
    
    if reg_user:
        message_text = f"✅ Добро пожаловать, {user.full_name}!\n\nВы успешно зарегистрированы!"
    else:
        message_text = "❌ Неверный код регистрации или ошибка сервера. Попробуйте снова."
    
    await state.clear()
    await message.answer(message_text)