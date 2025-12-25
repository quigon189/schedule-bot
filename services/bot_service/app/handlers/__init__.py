from aiogram import Router
from . import commands, user, admin, register
from app.callbacks.common_callbacks import echo_router

router = Router()

router.include_router(admin.admin_router)
router.include_router(register.register_router)
router.include_router(commands.router)
router.include_router(user.user_router)
router.include_router(echo_router)
