from aiogram import BaseMiddleware
from aiogram.types import Message, CallbackQuery
from app.services.user_service import user_service

class AdminMiddleware(BaseMiddleware):
    async def __call__(self, handler, event: Message | CallbackQuery, data: dict):
        user = await user_service.get_user(event.from_user.id)
        
        if user and "admin" in user.roles_list:
            return await handler(event, data)
        else:
            if isinstance(event, Message):
                await event.answer("❌ У вас недостаточно прав для этой команды")
            elif isinstance(event, CallbackQuery):
                await event.answer("❌ Недостаточно прав", show_alert=True)
            return