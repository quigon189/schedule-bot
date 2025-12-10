from aiogram import Router
from . import common_callbacks, user_callbacks

router = Router()

router.include_router(common_callbacks.router)
router.include_router(user_callbacks.router)