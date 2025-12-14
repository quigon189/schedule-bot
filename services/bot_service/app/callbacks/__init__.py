from aiogram import Router
from . import common_callbacks, user_callbacks, register_callback

router = Router()

router.include_router(common_callbacks.router)
router.include_router(user_callbacks.router)
router.include_router(register_callback.router)