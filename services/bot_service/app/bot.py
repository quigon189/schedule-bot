from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.fsm.storage.memory import MemoryStorage

from app.config import settings
from app.handlers import command_router, echo_router, admin_router, register_router
from app.middlewaries import CheckUserMiddleware

class TelegramBot:
    def __init__(self):
        self.bot = Bot(
            token=settings.TELEGRAM_BOT_TOKEN,
            default=DefaultBotProperties(parse_mode=ParseMode.HTML)
        )

        storage = MemoryStorage()
        self.dp = Dispatcher(storage=storage)

        self.dp.message.middleware(CheckUserMiddleware())

        self.dp.include_router(register_router)

        self.dp.include_router(admin_router)
        self.dp.include_router(command_router)
        self.dp.include_router(echo_router)

    async def start_polling(self):
        print("Start Telegram bot in polling mode...")
        await self.dp.start_polling(self.bot)

    async def stop(self):
        await self.bot.session.close()


telegram_bot = TelegramBot()
