import sys

import faker
import mysql.connector
from mysql.connector.cursor import MySQLCursorAbstract

fake = faker.Faker()


def create_table(cursor: MySQLCursorAbstract) -> None:
    create_stmt = """
    CREATE TABLE snippets (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(100) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
        expires DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL 7 DAY)
    );
    """
    try:
        cursor.execute("DROP TABLE IF EXISTS snippets;")
        cursor.execute(create_stmt)
    except mysql.connector.Error as err:
        print("Error while creating table: {} in line {}".format(err, sys.exc_info()[-1].tb_lineno))


def fill_table(cursor: MySQLCursorAbstract) -> None:
    try:
        cursor.execute("TRUNCATE TABLE snippets;")
        print("truncated")
        cursor.executemany(
            "INSERT INTO snippets (title, content) VALUES (%s, %s)",
            ((fake.sentence(nb_words=3), fake.text()) for _ in range(10)),
        )
        print("INSERTED")
    except mysql.connector.Error as err:
        print("Error while inserting data: {}".format(err))


if __name__ == "__main__":
    with (
        mysql.connector.connect(user="root", password="root", database="snippetbox") as conn,
        conn.cursor() as cursor,
    ):
        create_table(cursor)
        fill_table(cursor)
        conn.commit()
        # cursor.close()
        # conn.close()
