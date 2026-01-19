from typing import Literal, Optional, cast
import httpx
import logging
from cachetools import TTLCache
from app.config import settings
from app.models import CodeResponse, CreateUserRequest, ServiceResponse, UserResponse, CreateCodeRequest, Roles

logging.basicConfig(level=logging.DEBUG)


class AuthService:
    """
    Класс для взаимодействия с сервисом auth-service
    используем REST API
    с доступными методами можно ознакомится в swagger
    """

    def __init__(self):
        ttl = settings.AUTH_SERVICE_CACHE_TTL
        maxsize = settings.AUTH_SERVICE_CACHE_MAX
        self.base_url = settings.AUTH_SERVICE_URL
        self.timeout = settings.AUTH_SERVICE_TIMEOUT
        self._cache = TTLCache(maxsize=maxsize, ttl=ttl)

    async def create_registration_code(self, role: Roles, created_by: int,
                                       group_name: Optional[str] = None,
                                       max_uses: Optional[int] = None,
                                       expires: Optional[int] = None
                                       ) -> Optional[CodeResponse]:
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                payload = CreateCodeRequest(
                    role_name=cast(Roles, role),
                    created_by=created_by,
                    group_name=group_name,
                    max_uses=max_uses,
                    expiration=expires,
                )

                response = await client.post(
                    f"{self.base_url}/api/v1/code/create",
                    json=payload.model_dump(),
                    headers={"ContentType": "application/json"}
                )

                if response.status_code == 200:
                    auth_response = ServiceResponse(**response.json())
                    logging.debug(f"code responce: {auth_response}")
                    if auth_response.success:
                        return CodeResponse(**auth_response.data)

                return None

        except Exception as e:
            logging.debug(
                f"Error creating registration code in auth service: {e}")
            return None

    async def create_user(self, code: str, telegram_id: int,
                          username: Optional[str],
                          full_name: str) -> Optional[UserResponse]:
        """
        Создание запроса к auth-service на создание пользователя
        POST запрос к http://auth-service:8080/api/v1/users
        """
        try:
            # создаем асинхронного клиента для http запроса к auth-service
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                payload = CreateUserRequest(
                    code=code,
                    telegram_id=telegram_id,
                    username=username,
                    full_name=full_name
                )

                # асинхронно отправляем запрос и ждем ответ
                # (по умолчанию timeout = 30 sec)
                response = await client.post(
                    f"{self.base_url}/api/v1/users/register",
                    json=payload.model_dump(),
                    headers={"ContentType": "application/json"}
                )

                # TODO: необходимо переделать
                # проверяем ответ, если пользователь создан, то возвращаем его
                # иначе возвращаем None
                # (так же можно сделать лог пока через print)
                if response.status_code == 200:
                    auth_response = ServiceResponse(**response.json())
                    if auth_response.success:
                        data = UserResponse(**auth_response.data)
                        return data

                logging.debug(
                    f"Error: {response.status_code} {response.json()}")

                return None

        except Exception as e:
            logging.debug(f"Error creating user in auth service: {e}")
            return None

    async def get_user(self, telegram_id: int) -> Optional[UserResponse]:
        if telegram_id in self._cache:
            return self._cache[telegram_id]

        user = await self._get_user(telegram_id)

        if user:
            self._cache[telegram_id] = user

        return user

    async def invalidate_cache(self, telegram_id: Optional[int] = None):
        if telegram_id:
            self._cache.pop(telegram_id, None)
        else:
            self._cache.clear()

    async def _get_user(self, telegram_id: int) -> Optional[UserResponse]:
        """Получить пользователя с auth-service"""
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                response = await client.get(
                    f"{self.base_url}/api/v1/users/{telegram_id}"
                )

                if response.status_code == 200:
                    auth_response = ServiceResponse(**response.json())

                    if auth_response.success:
                        return UserResponse(**auth_response.data)

                return None

        except Exception as e:
            print(f"Error getting user from atuh service: {e}")
            return None

    # TODO: добавить остальные методы

    async def CheckIfRegistred():
        user = auth_service.get_user
        if user:
            IsRegistred = True
        else:
            return None
    
auth_service = AuthService()
