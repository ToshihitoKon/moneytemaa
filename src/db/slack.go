package db

func SlackInsertTransaction(price int, comment, slackUserId, slackChannelId, slackTimestamp string) error {
	db := GetDB()
	_, err := db.Query(
		"INSERT INTO transactions (price, comment, slack_user_id, slack_channel_id, slack_timestamp) VALUES (?, ?, ?, ?, ?)",
		price,
		comment,
		slackUserId,
		slackChannelId,
		slackTimestamp,
	)
	return err
}
