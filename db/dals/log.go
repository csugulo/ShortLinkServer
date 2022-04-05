package dals

import (
	"fmt"
	"sync"

	"github.com/csugulo/ShortLinkServer/consts"
	"github.com/csugulo/ShortLinkServer/db"

	log "github.com/sirupsen/logrus"
)

var writeMutex sync.Mutex

func AddLog(logType consts.LogType, url, urlID string, status consts.Status, statusMessage string) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()
	stmt, err := db.SqliteDB.Prepare(consts.InsertLogSQL)
	if err != nil {
		log.Errorf("prepare statement failed, err: %v", err)
		return err
	}
	res, err := stmt.Exec(logType, url, urlID, status, statusMessage)
	if err != nil {
		log.Errorf("execute query failed, err: %v", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Errorf("execute query failed, err: %v", err)
		return err
	} else {
		log.Infof("last insert id: %v", id)
	}
	return nil
}

type Statistic struct {
	LogType consts.LogType `json:"log_type"`
	Status  consts.Status  `json:"status"`
	Count   int64          `json:"count"`
}

func Statistics() ([]Statistic, error) {
	rows, err := db.SqliteDB.Query(consts.StatisticsSQL)
	if err != nil {
		log.Errorf("execute query failed, err: %v", err)
		return nil, err
	}
	var statistics []Statistic
	for rows.Next() {
		var statistic Statistic
		rows.Scan(
			&statistic.LogType,
			&statistic.Status,
			&statistic.Count,
		)
		statistics = append(statistics, statistic)
	}
	fmt.Println(statistics)
	return statistics, nil
}
