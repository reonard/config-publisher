package processor

import (
	"config-publisher/db"
	"fmt"
	"strconv"
	"time"
	"database/sql"
)

func GetDeviceInfo(deviceId int64) (error, DeviceInfo) {

	deviceInfo := DeviceInfo{}

	err := db.MySqlDB.QueryRow("select device_id, device_name, t_device.project as project from t_device "+
		"left join t_customer on project = t_customer.id where device_id = ? ", deviceId).Scan(&deviceInfo.DeviceId,
		&deviceInfo.DeviceName,
		&deviceInfo.Project)

	if err != nil {
		return err, deviceInfo
	}

	return nil, deviceInfo

}

func InsertAlarmItem(data *MonitorData, metricKey string) {

	unixTime, _ := strconv.Atoi(data.TimeStamp)

	err := db.ExecuteSQL(
		"INSERT INTO t_alarm_item (`device_id`, `timestamp`, `alarm_item`, `alarm_value`) VALUES (?,?,?,?)",
		data.DeviceId, time.Unix(int64(unixTime/1000), 0), metricKey, data.Data[0][metricKey])

	if err != nil {
		fmt.Println(err)
	}
}


func GetLatestPushData(minutesWaitForRetry int64) (*sql.Rows, error){

	rows, err := db.MySqlDB.Query(
		"SELECT ts.device_id,ts.id,tm.action_type,tc.configs " +
			"FROM t_issue_status ts " +
			"LEFT JOIN t_issue_msg tm ON ts.issue_msg_id = tm.id " +
			"LEFT JOIN t_device_config tc ON tm.config_id = tc.id " +
			"WHERE ts.issue_status = ? OR " +
			"(ts.issue_status != ? AND ts.status_time < DATE_SUB(CURRENT_TIMESTAMP, INTERVAL ? MINUTE))", STATUS_PENDING, STATUS_ACK, minutesWaitForRetry)

	return rows, err
}

func UpdateIssueStatus(statusId int64, msgStatus int8){

	err := db.ExecuteSQL("UPDATE t_issue_status " +
		"SET issue_status = ?, status_time = current_timestamp()" +
		"WHERE id = ?", msgStatus, statusId)

	if err != nil {
		fmt.Println(err)
	}

}