package sequence

import (
	"database/sql"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 基于mysql实现取号器
const sqlReplaceStub = `REPLACE INTO sequence (stub) VALUES ('a')`

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

// Next 取下一个号
func (m *MySQL) Next() (u uint64, err error) {
	// 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceStub)
	if err != nil {
		logx.Errorw("conn.Prepare failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return 0, err
	}
	defer func(stmt sqlx.StmtSession) {
		err := stmt.Close()
		if err != nil {
			logx.Errorw("stmt.Close failed", logx.LogField{
				Key:   "err",
				Value: err.Error(),
			})
		}
	}(stmt)

	// 执行
	var res sql.Result
	res, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return 0, err
	}

	// 拿到刚插入的主键id
	var num int64
	num, err = res.LastInsertId()
	if err != nil {
		logx.Errorw("res.LastInsertId failed", logx.LogField{
			Key:   "err",
			Value: err.Error(),
		})
		return 0, err
	}
	return uint64(num), nil
}
