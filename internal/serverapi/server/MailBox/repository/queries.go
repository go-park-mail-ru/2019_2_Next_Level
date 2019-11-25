package repository

const(
	queryGetEmailByCode = `SELECT sender, email AS "receivers", time, subject, body from Message
							JOIN Receiver ON Message.id=Receiver.mailId
							WHERE Message.id=$1`

	queryGetEmailList = `SELECT Message.id, sender, email AS "receivers", time, subject, body, isread from Message
						JOIN Receiver ON Message.id=Receiver.mailId
						WHERE %s=$1 ORDER BY time LIMIT $2 OFFSET $3;`

	queryGetMessagesCount = `SELECT COUNT(Message.id) from Message JOIN Receiver ON Message.id=Receiver.mailId
							WHERE Receiver.email=$1`

	queryMarkMessage =  `UPDATE Message SET %s WHERE id=$1`

	queryWriteMessage = `INSERT INTO Message (sender, subject, body, direction, folder) VALUES
							($1, $2, $3, $4, $5)
						RETURNING id;`

	queryInflateReceivers = `INSERT INTO Receiver (mailId, email) VALUES`
)
