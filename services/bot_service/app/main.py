import asyncio
from contextlib import asynccontextmanager

from fastapi import FastAPI
import uvicorn

from app.bot import telegram_bot
from app.web.routes import router
from app.config import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    bot_task = asyncio.create_task(telegram_bot.start_polling())

    yield

    print("Shutdown telegram bot...")
    bot_task.cancel()
    try:
        await bot_task
    except asyncio.CancelledError:
        pass
    await telegram_bot.stop()

app = FastAPI(
    title="Telegram Bot Service",
    discription="Микросервис для работы с Telegram ботом",
    version="1.0.0",
    lifespan=lifespan
)

app.include_router(router)

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=8000,
        reload=settings.DEBUG,
        log_level="debug"
    )
