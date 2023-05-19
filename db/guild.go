package db

import (
	"database/sql"

	tempest "github.com/Amatsagu/Tempest"
)

func SelectAllValidChannels(tx *sql.Tx) ([]tempest.Snowflake, error) {
	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.Query("SELECT channel_id FROM guild WHERE blocked = false;")
	} else {
		rows, err = Conn.Query("SELECT channel_id FROM guild WHERE blocked = false;")
	}

	if err != nil {
		return nil, err
	}

	result := make([]tempest.Snowflake, 0)
	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}

		var channelID tempest.Snowflake
		err := rows.Scan(&channelID)
		if err != nil {
			return nil, err
		}

		result = append(result, channelID)
	}

	return result, nil // < moved []card{} to heap
}
