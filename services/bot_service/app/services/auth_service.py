from typing import Optional
import httpx
import logging
from app.config import settings
from app.models import CreateUserRequest, AuthResponse, UserResponse

logging.basicConfig(level=logging.DEBUG)


class AuthService:
    """
    Класс для взаимодействия с сервисом auth-service
    используем REST API
    с доступными методами можно ознакомится в swagger
    """

    def __init__(self):
        self.base_url = settings.AUTH_SERVICE_URL
        self.timeout = settings.AUTH_SERVICE_TIMEOUT

    async def create_user(self, telegram_id: int, username: Optional[str],
                          full_name: str) -> Optional[UserResponse]:
        """
        Создание запроса к auth-service на создание пользователя
        POST запрос к http://auth-service:8080/api/v1/users
        """
        try:
            # создаем асинхронного клиента для http запроса к auth-service
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                payload = CreateUserRequest(
                    telegram_id=telegram_id,
                    username=username,
                    full_name=full_name
                )

                # асинхронно отправляем запрос и ждем ответ
                # (по умолчанию timeout = 30 sec)
                response = await client.post(
                    f"{self.base_url}/api/v1/users",
                    json=payload.model_dump(),
                    headers={"ContentType": "application/json"}
                )

                logging.debug(f"response data: {response.json()}")

                # TODO: необходимо переделать
                # проверяем ответ, если пользователь создан, то возвращаем его
                # иначе возвращаем None
                # (так же можно сделать лог пока через print)
                if response.status_code == 200:
                    auth_response = AuthResponse(**response.json())
                    if auth_response.success:
                        return auth_response.data

                logging.debug(
                    f"Error: {response.status_code} {response.json()}")

                return None

        except Exception as e:
            logging.debug(f"Error creating user in auth service: {e}")
            return None

    async def get_user(self, telegram_id: int) -> Optional[UserResponse]:
        """Получить пользователя с auth-service"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.get(
                    f"{self.base_url}/api/v1/users/{telegram_id}"
                )

                if response.status_code == 200:
                    auth_response = AuthResponse(**response.json())

                    if auth_response.success:
                        return auth_response.data

                return None

        except Exception as e:
            print(f"Error getting user from atuh service: {e}")
            return None

    # TODO: добавить остальные методы


auth_service = AuthService()
