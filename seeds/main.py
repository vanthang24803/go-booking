from faker import Faker
import psycopg2 as pg
import logging
from datetime import datetime
import random

logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s"
)

fake = Faker(["vi_VN", "en_US", "ja_JP"])


def connect_db():
    try:
        conn = pg.connect(
            database="go-ecommerce",
            user="postgres",
            password="Thang@240803",
            host="localhost",
            port=5432,
        )
        logging.info("Connected to PostgreSQL successfully! ✅")
        return conn
    except Exception as e:
        logging.error(f"Error connecting to PostgreSQL: {e}")
        return None


def generate_random_users(n):
    users = []
    default_password = "$2a$10$NdXZuedsz.QT01XWCJtmwe5RGM5ZE5q9xCsAV8J61SzYaZt4A72Xa"

    for _ in range(n):
        first_name = fake.first_name()
        surname = fake.last_name()
        username = f"{first_name.lower()}.{surname.lower()}{random.randint(1, 1000)}"
        email = f"{username}@{fake.free_email_domain()}"
        hash_password = default_password
        avatar = fake.image_url()
        created_at = datetime.now()
        updated_at = created_at

        users.append(
            (
                username,
                email,
                hash_password,
                first_name,
                surname,
                avatar,
                created_at,
                updated_at,
            )
        )
    return users


def seed_user_roles(conn, start_user_id, end_user_id, role_id):
    """Insert user roles with specific user_id and role_id."""
    try:
        cursor = conn.cursor()
        query = """
        INSERT INTO user_roles (user_id, role_id)
        VALUES (%s, %s)
        """
        data = [(user_id, role_id) for user_id in range(start_user_id, end_user_id + 1)]
        cursor.executemany(query, data)
        conn.commit()
        logging.info(
            f"Seeded user_roles for user_id {start_user_id} -> {end_user_id} with role_id {role_id} successfully! ✅"
        )
    except Exception as e:
        logging.error(f"Error seeding user_roles: {e}")
        conn.rollback()


def seed_users(conn, users):
    try:
        cursor = conn.cursor()
        query = """
        INSERT INTO users (username, email, hash_password, first_name, surname, avatar, created_at, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
        """
        cursor.executemany(query, users)
        conn.commit()
        logging.info(f"Seeded {len(users)} users successfully! ✅")
    except Exception as e:
        logging.error(f"Error seeding users: {e}")
        conn.rollback()


def main():
    conn = connect_db()
    if conn:
        users = generate_random_users(180)
        seed_users(conn, users)

        conn.close()


if __name__ == "__main__":
    main()
