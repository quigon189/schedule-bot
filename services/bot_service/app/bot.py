from typing import List, Optional
from aiogram import Bot, Dispatcher
from aiogram.client.default import DefaultBotProperties
from aiogram.enums import ParseMode
from aiogram.fsm.storage.memory import MemoryStorage
from aiogram.utils.media_group import MediaGroupBuilder

from app.config import settings
from app.handlers import router
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

        self.dp.include_router(router)

    async def start_polling(self):
        print("Start Telegram bot in polling mode...")
        await self.dp.start_polling(self.bot)

    async def send_message(self, chat_id: int,
                           message: str = "",
                           photo_urls: Optional[List[str]] = None):
        if photo_urls:
            album_builder = MediaGroupBuilder(caption=message)
            for photo_url in photo_urls:
                album_builder.add_photo(media=photo_url)

            await self.bot.send_media_group(
                chat_id=chat_id,
                media=album_builder.build()
            )
            return

        await self.bot.send_message(chat_id=chat_id, text=message)

    async def stop(self):
        await self.bot.session.close()


telegram_bot = TelegramBot()
