package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/genshinsim/gcsim/backend/pkg/services/submission"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	DBPath string
}

type Store struct {
	cfg Config
	Log *zap.SugaredLogger
	db  *sql.DB
}

func New(cfg Config, cust ...func(*Store) error) (*Store, error) {
	s := &Store{
		cfg: cfg,
	}

	var err error

	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	err = s.createOrConnectToDB()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) createOrConnectToDB() error {
	//if path does not exist, create new + initialize schema
	//otherwise assume schema is correct
	//check if file exists, if not create new and update schema
	create := false
	if _, err := os.Stat(s.cfg.DBPath); os.IsNotExist(err) {
		// path/to/whatever exists
		create = true
	} else if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", s.cfg.DBPath)
	if err != nil {
		return fmt.Errorf("error opening db: %v", err)
	}
	if create {
		s.Log.Infow("db not found, building new")
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error setting up new schema: %v", err)
		}
	}
	s.db = db

	return nil

}

const schema = `
create table data (
	id text not null primary key
	, config text not null default ""
	, description text not null default ""
	, submitter text not null default ""
	, preview text not null default ""
	, to_compute integer not null default 1
);
`

func (s *Store) Get(id string) (*submission.Submission, error) {
	s.Log.Infow("submission (sqlite): get request received", "id", id)

	const sqlstr = `SELECT id, config, description, submitter, preview FROM data where id = ?`

	var res submission.Submission

	err := s.db.QueryRow(sqlstr, id).Scan(
		&res.Id,
		&res.Config,
		&res.Description,
		&res.Submitter,
		&res.Preview,
	)

	if err != nil {
		s.Log.Infow("submission (sqlite): get request error", "id", id, "err", err)
		if err == sql.ErrNoRows {
			return nil, status.Error(codes.NotFound, "id not found")
		}
		return nil, status.Error(codes.Internal, "interal server error")
	}

	return &res, nil
}

func (s *Store) Set(req *submission.Submission) error {
	s.Log.Infow("submission (sqlite): set request received", "req", req.String())

	//sanity check if id exists first; strictly speaking not required since server already does this check
	var err error
	err = s.idCheck(req.GetId())
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`update data set 
	config = $1, description = $2, to_compute = $3 where id = $4`,
		req.GetConfig(),
		req.GetDescription(),
		1, //force recompute
		req.GetId(),
	)

	if err != nil {
		s.Log.Infow("submission (sqlite): set request error", "err", err, "req", req.String())
		return status.Error(codes.Internal, "internal server error")
	}

	return nil
}

func (s *Store) Delete(id string) error {
	s.Log.Infow("submission (sqlite): delete request received", "id", id)

	//sanity check if id exists first; strictly speaking not required since server already does this check
	var err error
	err = s.idCheck(id)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`delete from data where id = ?`, id)

	if err != nil {
		s.Log.Infow("submission (sqlite): delete request error", "err", err, "id", id)
		return status.Error(codes.Internal, "internal server error")
	}

	return nil
}

func (s *Store) New(req *submission.Submission) (string, error) {
	s.Log.Infow("submission (sqlite): new entry request received", "req", req.String())
	id, err := nanoid.New()
	if err != nil {
		s.Log.Infow("submission (sqlite): new entry request unexpected err generating unique id", "err", err, "req", req.String())
	}
	//sanity checks
	if req.GetConfig() == "" {
		return "", status.Error(codes.InvalidArgument, "config cannot be blank")
	}
	if req.GetSubmitter() == "" {
		return "", status.Error(codes.InvalidArgument, "submitter cannot be blank")
	}

	const sqlstr = `insert into data (
		id, config, description, submitter 
	) values (
		$1, $2, $3, $4 
	)
	`

	if _, err := s.db.Exec(
		sqlstr,
		id,
		req.GetConfig(),
		req.GetDescription(),
		req.GetSubmitter(),
	); err != nil {
		s.Log.Infow("submission (sqlite): new entry request unexpected err", "err", err, "req", req.String())
		return "", status.Error(codes.Internal, "internal server error")
	}
	s.Log.Infow("submission (sqlite): new entry added", "id", id, "req", req.String())

	return id, nil
}

func (s *Store) List(filter string) ([]*submission.Submission, error) {
	s.Log.Infow("submission (sqlite): list request received", "filter", filter)

	const sqlstr = `SELECT id, config, description, submitter, preview FROM data`

	rows, err := s.db.Query(sqlstr)
	if err != nil {
		s.Log.Infow("submission (sqlite): list request error querying rows", "err", err)
		return nil, status.Error(codes.Internal, "interal server error")
	}

	var res []*submission.Submission
	filters := strings.Split(filter, ",")
	hasFilter := len(filters) > 0

	defer rows.Close()
nextRow:
	for rows.Next() {
		var t submission.Submission
		err := rows.Scan(
			&t.Id,
			&t.Config,
			&t.Description,
			&t.Submitter,
			&t.Preview,
		)

		if err != nil {
			s.Log.Infow("submission (sqlite): list request error scanning row", "err", err)
			return nil, status.Error(codes.Internal, "interal server error")
		}

		if hasFilter {
			// if submitter not in any of the filter skip
			found := false
			for _, v := range filters {
				if v == t.Submitter {
					found = true
					break
				}
			}
			if !found {
				continue nextRow
			}
		}
		res = append(res, &t)
	}
	err = rows.Err()

	if err != nil {
		s.Log.Infow("submission (sqlite): list request row contain error", "err", err)
		return nil, status.Error(codes.Internal, "interal server error")
	}

	return res, nil
}

func (s *Store) idCheck(id string) error {
	s.Log.Infow("submission (sqlite): check for id", "id", id)
	if id == "" {
		return status.Error(codes.InvalidArgument, "invalid id")
	}
	var count int
	err := s.db.QueryRow(`select count(id) from data where id = ?`, id).Scan(&count)
	if err != nil {
		s.Log.Infow("submission (sqlite): id check error", "err", err, "id", id)
		return status.Error(codes.Internal, "internal server error")
	}
	if count == 0 {
		s.Log.Infow("submission (sqlite): id check error", "err", "id does not exist", "id", id)
		return status.Error(codes.NotFound, "id not found")
	}

	return nil
}
