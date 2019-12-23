package repository

const(
	queryGetEmailByCode = `SELECT sender, email AS "receivers", time, subject, body, isRead from Message
							JOIN Receiver ON Message.id=Receiver.mailId
							WHERE Message.id=$1`

	queryGetEmailList = `SELECT Message.id, sender, email AS "receivers", time, subject, body, isread from Message
						JOIN Receiver ON Message.id=Receiver.mailId
						WHERE Message.owner=$1 AND folder=$2 AND Message.id<$4 ORDER BY id DESC LIMIT $3;`

	queryGetMessagesCount = `SELECT COUNT(Message.id) from Message JOIN Receiver ON Message.id=Receiver.mailId
							WHERE Receiver.email=$1`

	queryMarkMessage =  `UPDATE Message SET %s WHERE id=$1`

	queryWriteMessage = `INSERT INTO Message (sender, subject, body, direction, folder, owner) VALUES
							($1, $2, $3, $4, $5, $6)
						RETURNING id;`

	queryInflateReceivers = `INSERT INTO Receiver (mailId, email) VALUES`
)
