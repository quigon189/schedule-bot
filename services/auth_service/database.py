from models import UserInDB, Role
from passlib.context import CryptContext

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")


def verify_password(plain_password, hashed_password) -> bool:
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(password) -> str:
    return pwd_context.hash(password)


def get_user(db, username: str):
    pass


def authenticate_user(username: str, password: str):
    user = get_user("db", username)
    if not user:
        return False
    if not verify_password(password, user.hashed_password):
        return False
    return user
