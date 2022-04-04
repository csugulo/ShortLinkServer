package consts

type LogType int8

const (
	Create   LogType = 1
	Redirect LogType = 2
)

func (logType LogType) String() string {
	switch logType {
	case Create:
		return "Create"
	case Redirect:
		return "Redirect"
	default:
		return "Unkown"
	}
}

type Status int8

const (
	Success Status = 0
	Failed  Status = 1
)

func (status Status) String() string {
	switch status {
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	default:
		return "Unkown"
	}
}

const CreateLogTable = `
create table if not exists log(
	id integer primary key autoincrement,
	log_type varchar(255),
	url text,
	url_id varchar(255),
	status integer,
	status_message varchar(255),
	create_time text default current_timestamp
)
`

const InsertLogSQL = `
insert into log(log_type, url, url_id, status, status_message) values(?, ?, ?, ?, ?)
`

const StatisticsSQL = `
select log_type, status, count(1) from log group by log_type, status
`
