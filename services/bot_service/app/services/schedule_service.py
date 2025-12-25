from datetime import datetime
import logging
from typing import List, Optional

import httpx
from app.config import settings
from app.models import GroupScheduleResponse, GroupScheduleRequest, ScheduleChangesResponse, ServiceResponse

logging.basicConfig(level=logging.DEBUG)


class ScheduleService:
    """
    Класс для взаимодействия с schedule-service
    """

    def __init__(self):
        self.base_url = f"{settings.SCHEDULE_SERVICE_URL}/api/v1"
        self.timeout = settings.SCHEDULE_SERVICE_TIMEOUT

    async def get_group_schedule(self, group_name: str, academic_year: str,
                                 half_year: int) -> Optional[List[GroupScheduleResponse]]:
        """
        Возвращает основное расписание группы на указаный учебный период
        """
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                payload = GroupScheduleRequest(
                    academic_year=academic_year,
                    half_year=half_year,
                    group_name=group_name
                )

                response = await client.get(
                    f"{self.base_url}/group_schedules",
                    params=payload.model_dump(),
                    headers={"ContentType": "application/json"}
                )

                if response.status_code == 200:
                    schedule_response = ServiceResponse(**response.json())
                    if schedule_response.success:
                        gr_schedules = [GroupScheduleResponse(
                            **sh) for sh in schedule_response.data]
                        return gr_schedules

                logging.debug(
                    f"Error: {response.status_code} {response.json()}")

                return None

        except Exception as e:
            logging.debug(f"Error getting group schedule: {e}")
            return None

    async def get_schedule_changes(self, date: Optional[datetime] = None) -> Optional[ScheduleChangesResponse]:
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                params = {}
                if date:
                    params['date'] = date.strftime('%Y-%m-%d')

                response = await client.get(
                    f"{self.base_url}/changes",
                    params=params,
                    headers={"ContentType": "application/json"}
                )

                logging.debug("response: %+v")

                if response.status_code == 200:
                    changes_response = ServiceResponse(**response.json())
                    if changes_response.success:
                        return ScheduleChangesResponse(**changes_response.data)

                logging.debug(
                    f"Error: {response.status_code} {response.json()}")

                return None

        except Exception as e:
            logging.debug(f"Error getting schedule changes: {e}")
            return None



schedule_service = ScheduleService()
