package repository

const (
	CreateTableUsers = `
 	CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        birthday DATE NOT NULL
    );`

	CreateTableSubscriptions = `
	CREATE TABLE IF NOT EXISTS subscriptions (
        subscriber_id INTEGER REFERENCES users(id),
        birthday_user_id INTEGER REFERENCES users(id),
        PRIMARY KEY (subscriber_id, birthday_user_id)
    );`

	addUser = `INSERT INTO users (email, name, birthday) VALUES ($1, $2, $3);`

	subscribe = `INSERT INTO subscriptions (subscriber_id, birthday_user_id) VALUES ($1, $2);`

	unsubscribe = `DELETE FROM subscriptions WHERE subscriber_id = $1 AND birthday_user_id = $2;`

	getSubscribersForTodayBirthdays = `
	SELECT sub.email, bday.name, bday.email
    FROM subscriptions s
    JOIN users sub ON s.subscriber_id = sub.id
    JOIN users bday ON s.birthday_user_id = bday.id
    WHERE TO_CHAR(bday.birthday, 'MM-DD') = TO_CHAR(current_date, 'MM-DD');`
)
