from aiogram import F, Router
from aiogram.filters import Command
from aiogram.types import CallbackQuery, InlineKeyboardButton, InlineKeyboardMarkup, Message
from app.middlewares.admin import AdminMiddleware
from app.services.auth_service import auth_service

admin_router = Router()
admin_router.message.middleware(AdminMiddleware())

groups = ["СА-501", "СА-502"]

def get_code_type_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [
                InlineKeyboardButton(text="👨‍🎓 Студент", callback_data="create_code:student"),
                InlineKeyboardButton(text="👨‍🏫 Преподаватель", callback_data="create_code:teacher"),
                InlineKeyboardButton(text="👑 Менеджер", callback_data="create_code:manager")
            ]
        ]
    )
    return keyboard

def get_groups_keyboard():
    keyboard = InlineKeyboardMarkup(
        inline_keyboard=[
            [InlineKeyboardButton(text=group, callback_data=f"group:{group}") for group in groups]
        ]
    )
    return keyboard

@admin_router.message(Command("code"))
async def create_code(message: Message):
    await message.answer("Выберите роль для кода:", reply_markup=get_code_type_keyboard())

@admin_router.callback_query(F.data.startswith('create_code:'))
async def process_code_role(callback_query: CallbackQuery):
    role = callback_query.data.split(':')[1]
    
    if role == 'student':
        await callback_query.message.edit_text(
            'Выберите группу:',
            reply_markup=get_groups_keyboard()
        )
    else:
        await create_final_code(callback_query, role, None)
    
    await callback_query.answer()

@admin_router.callback_query(F.data.startswith('group:'))
async def process_group_selection(callback_query: CallbackQuery):
    group = callback_query.data.split(':')[1]
    role = 'student'
    await create_final_code(callback_query, role, group)

async def create_final_code(callback_query: CallbackQuery, role: str, group: str = None):
    code = await auth_service.create_registration_code(
        role=role,
        group_name=group,
        max_uses=10,
        created_by=callback_query.from_user.id,
        expires=7 * 24 * 60 * 60  # 7 дней
    )
    
    if code:
        result_text = (
            f"✅ Код для {role} создан!\n\n"
            f"🔑 Код: `{code.code}`\n"
            f"👥 Группа: {group or 'любая'}\n"
            f"🔢 Использований: {code.max_uses}\n"
            f"⏰ Действует до: {code.expires_at.strftime('%d.%m.%Y %H:%M')}"
        )
    else:
        result_text = "❌ Ошибка при создании кода"
    
    await callback_query.message.edit_text(result_text)
    await callback_query.answer()