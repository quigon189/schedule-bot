from aiogram import F, Router
from aiogram.types import CallbackQuery, Message
from aiogram.fsm.state import State, StatesGroup
from aiogram.fsm.context import FSMContext
from app.models import TelegramUser
from app.services.user_service import user_service
from app.keyboards.user_keyboards import get_main_menu_keyboard

register_router = Router()


class RegisterState(StatesGroup):
    waiting_for_code = State()


@register_router.message(RegisterState.waiting_for_code)
async def process_register_code(message: Message, state: FSMContext):
    code = message.text.strip()

    if not code:
        await message.answer("❌ Код не может быть пустым. Попробуйте снова:")
        return

    # Создаем объект пользователя Telegram
    user = TelegramUser(
        id=message.from_user.id,
        username=message.from_user.username,
        first_name=message.from_user.first_name,
        last_name=message.from_user.last_name
    )

    # Пытаемся зарегистрировать пользователя через auth_service
    reg_user = await user_service.register_user(code, user)

    if reg_user:
        # Очищаем состояние
        await state.clear()

        # Очищаем кэш для этого пользователя
        from app.services.auth_service import auth_service
        await auth_service.invalidate_cache(message.from_user.id)

        # Приветствуем пользователя
        await message.answer(
            f"✅ Регистрация успешна!\n\n"
            f"Добро пожаловать, {reg_user.full_name}!\n\n"
            f"🎓 Ваши данные:\n"
            f"• Роли: {', '.join(reg_user.roles_list) if reg_user.roles_list else 'не назначены'}\n"
            f"• Группа: {reg_user.group or 'не указана'}\n"
            f"• Статус: {'✅ Активен' if reg_user.is_active else '❌ Неактивен'}\n\n"
            f"Теперь вы можете использовать все функции бота!",
            reply_markup=get_main_menu_keyboard()
        )
    else:
        await message.answer(
            "❌ Регистрация не удалась!\n\n"
            "Возможные причины:\n"
            "1. Неверный код регистрации\n"
            "2. Код истек или использован максимальное количество раз\n"
            "3. Проблемы с сервером авторизации\n\n"
            "Проверьте код и попробуйте снова.\n"
            "Для повторной попытки нажмите 'Зарегистрироваться'.",
            reply_markup=get_main_menu_keyboard()
        )
        await state.clear()
