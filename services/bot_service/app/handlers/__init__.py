from .commands import command_router
from .echo import echo_router
from .admin import admin_router
from .register import register_router

__all__ = ['command_router', 'echo_router', 'admin_router', 'register_router']
