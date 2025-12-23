from datetime import datetime
from fastapi import APIRouter

from app.bot import telegram_bot
from app.models import HealthCheck, TgSendRequest


router = APIRouter()


@router.get("/health", response_model=HealthCheck)
async def health_check():
    return HealthCheck(
        status="heathy",
        service="bot-service",
        timestampt=datetime.now()
    )


@router.post("/send_message")
async def send_message(req: TgSendRequest):
    await telegram_bot.send_message(
        chat_id=req.chat_id,
        message=req.message,
        photo_urls=req.photo_urls
    )
