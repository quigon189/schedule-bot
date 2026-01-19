from aiogram import F, Router
from aiogram.filters import Command
from aiogram.types import CallbackQuery, InlineKeyboardButton, InlineKeyboardMarkup, Message

from app.middlewaries.admin import AdminMiddleware
from app.services.auth_service import auth_service


admin_router = Router()
admin_router.message.middleware(AdminMiddleware())

groups = ["СА-501", "СА-502"]


def get_code_type_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(
                    text="Менеджер", callback_data="create_code:manager"),
                InlineKeyboardButton(
                    text="Преподаватель", callback_data="create_code:teacher"),
                InlineKeyboardButton(
                    text="Студент", callback_data="create_code:student")

            ]
        ]
    )
    return keyboard


def get_groups_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(
                    text=f"{group}", callback_data=f"create_code:student:{group}")
                for group in groups
            ]
        ]
    )
    return keyboard


@admin_router.message(Command("code"))
async def create_code(message: Message):
    await message.answer("Выберите роль:", reply_markup=get_code_type_keyboard())


@admin_router.callback_query(F.data.startswith('create_code:'))
async def create_code_callback(callback_query: CallbackQuery):
    data_parts = callback_query.data.split(':')
    role = data_parts[1]

    if role == 'student':
        if len(data_parts) < 3:
            await callback_query.message.edit_text(
                'Укажите название группы:',
                reply_markup=get_groups_keyboard()
            )
            await callback_query.answer()
            return

        group = data_parts[2]

        code = await auth_service.create_registration_code(
            role=role,
            group_name=group,
            max_uses=30,
            created_by=callback_query.from_user.id
        )
    else:
        code = await auth_service.create_registration_code(
            role=role,
            created_by=callback_query.from_user.id
        )

    if code:
        result_text = f"""
Код доступа для {role} создан
Код: {code.code}
Количество использований: {code.max_uses}
Действителен до: {code.expires_at.strftime("%Y-%y-%d %H:%M:%S")}
        """
    else:
        result_text = "При создании кода возникла ошибка"

    await callback_query.message.edit_text(result_text)
    await callback_query.answer()
