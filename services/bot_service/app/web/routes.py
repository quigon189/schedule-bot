from datetime import datetime
from fastapi import APIRouter


from app.models import HealthCheck


router = APIRouter()


@router.get("/health", response_model=HealthCheck)
async def health_check():
    return HealthCheck(
        status="heathy",
        service="bot-service",
        timestampt=datetime.now()
    )
