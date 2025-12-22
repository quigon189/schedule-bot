from datetime import datetime
import logging
from typing import Optional

from app.models import AiResponse
import httpx


class AiService:
    def __init__(self):
        self.base_url = "http://n8n:5678/webhook"
        self.timeout = 600

    async def proccess(self, message: str) -> Optional[AiResponse]:
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                # пока что вышиты статические параметры для примера
                payload = {
                    "user_message": {
                        "message": message,
                        "current_date": datetime.now().strftime("%Y-%m-%d"),
                        "current_academic_year": "2025/2026",
                        "current_half_year": 1,
                    }
                }

                response = await client.post(
                    f"{self.base_url}/ai",
                    json=payload,
                    headers={"ContentType": "appliction/json"}
                )

                logging.debug(f"ai response text: {response.text}")
                logging.debug(f"ai response headers: {response.headers}")

                if response.status_code == 200:
                    return AiResponse(**response.json())

            return None
        except Exception as e:
            logging.debug(f"Error with proccess ai agent: {e}")
            return None


n8n_service = AiService()
