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
