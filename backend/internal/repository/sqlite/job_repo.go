package sqlite

import (
	"context"
	"database/sql"

	"matasidik/backend/internal/domain"
)

func (s *Store) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	var out domain.Job
	err := s.db.GetContext(ctx, &out, `SELECT id, type, payload_json, dedupe_key, status, attempts,
		max_attempts, run_after, locked_by, locked_at, heartbeat_at, last_error, progress_json,
		created_at, finished_at FROM job WHERE id = ?`, id)
	if err == sql.ErrNoRows {
		return nil, domain.NotFound("job")
	}
	return &out, err
}

func (s *Store) ClaimJob(ctx context.Context, workerID, now string) (*domain.Job, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var job domain.Job
	err = tx.GetContext(ctx, &job, `SELECT id, type, payload_json, dedupe_key, status, attempts,
		max_attempts, run_after, locked_by, locked_at, heartbeat_at, last_error, progress_json,
		created_at, finished_at FROM job
		WHERE status = 'queued' AND run_after <= ?
		ORDER BY created_at LIMIT 1`, now)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE job SET status = 'running', locked_by = ?,
		locked_at = ?, heartbeat_at = ?, attempts = attempts + 1 WHERE id = ?`,
		workerID, now, now, job.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	job.Status = domain.JobRunning
	job.Attempts++
	job.LockedBy = &workerID
	job.LockedAt = &now
	job.HeartbeatAt = &now
	return &job, nil
}

func (s *Store) UpdateJobProgress(ctx context.Context, id, progressJSON string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE job SET progress_json = ?, heartbeat_at = ? WHERE id = ?`,
		progressJSON, nowRFC3339(), id)
	return err
}

func (s *Store) CompleteJob(ctx context.Context, id, progressJSON, finishedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE job SET status = 'done', progress_json = ?,
		finished_at = ?, heartbeat_at = ? WHERE id = ?`, progressJSON, finishedAt, finishedAt, id)
	return err
}

func (s *Store) FailJob(ctx context.Context, id, errMsg, finishedAt string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE job
		SET status = CASE WHEN attempts >= max_attempts THEN 'dead' ELSE 'failed' END,
		    last_error = ?, finished_at = ?, heartbeat_at = ? WHERE id = ?`,
		errMsg, finishedAt, finishedAt, id)
	return err
}

func (s *Store) RequeueStaleJobs(ctx context.Context, before string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE job SET status = 'queued', locked_by = NULL,
		locked_at = NULL WHERE status = 'running' AND heartbeat_at < ?`, before)
	return err
}
