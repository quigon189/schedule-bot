from aiogram import Router
from . import commands, user, echo

router = Router()

router.include_router(commands.router)
router.include_router(user.router)
router.include_router(echo.router)