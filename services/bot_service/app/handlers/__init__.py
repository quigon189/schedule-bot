from aiogram import Router
from . import commands, user, echo, admin

router = Router()

router.include_router(commands.router)
router.include_router(user.router)
router.include_router(echo.router)

# admin роутер должен подключаться отдельно с middleware