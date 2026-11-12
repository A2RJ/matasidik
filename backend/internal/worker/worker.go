package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"matasidik/backend/internal/domain"
	"matasidik/backend/internal/repository"
)

type Worker struct {
	id     string
	subs   repository.SubmissionRepository
	jobs   repository.JobRepository
	master repository.MasterRepository
	class  *Classifier
}

func New(id string, subs repository.SubmissionRepository, jobs repository.JobRepository, master repository.MasterRepository) *Worker {
	return &Worker{id: id, subs: subs, jobs: jobs, master: master, class: NewClassifier(master)}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()
	_ = w.jobs.RequeueStaleJobs(ctx, time.Now().Add(-5*time.Minute).UTC().Format(time.RFC3339))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if _, err := w.ProcessOnce(ctx); err != nil {
				slog.Error("worker", "error", err)
			}
		}
	}
}

func (w *Worker) ProcessOnce(ctx context.Context) (bool, error) {
	job, err := w.jobs.ClaimJob(ctx, w.id, time.Now().UTC().Format(time.RFC3339))
	if err != nil || job == nil {
		return false, err
	}
	if err := w.process(ctx, job); err != nil {
		_ = w.jobs.FailJob(ctx, job.ID, err.Error(), time.Now().UTC().Format(time.RFC3339))
		return true, err
	}
	return true, nil
}

type jobPayload struct {
	SubmissionID string `json:"submission_id"`
	FileID       string `json:"file_id"`
}

func (w *Worker) process(ctx context.Context, job *domain.Job) error {
	var payload jobPayload
	if err := json.Unmarshal([]byte(job.PayloadJSON), &payload); err != nil {
		return fmt.Errorf("payload job tidak valid: %w", err)
	}
	sub, err := w.subs.GetSubmission(ctx, payload.SubmissionID)
	if err != nil {
		return err
	}
	file, err := w.subs.GetFileObject(ctx, sub.FileID)
	if err != nil {
		return err
	}
	if err := w.progress(ctx, job.ID, "validasi", 15, nil); err != nil {
		return err
	}
	raw, err := ParseFile(file.Path)
	if err != nil {
		return fmt.Errorf("parse file: %w", err)
	}
	if err := w.progress(ctx, job.ID, "klasifikasi", 45, map[string]any{"rows": len(raw)}); err != nil {
		return err
	}
	classified, err := w.class.Classify(ctx, sub, raw)
	if err != nil {
		return err
	}
	matched := 0
	for _, r := range classified.Rows {
		if r.Classification == domain.ClassMatched {
			matched++
		}
	}
	if err := w.progress(ctx, job.ID, "simpan", 80, map[string]any{"rows": len(classified.Rows), "matched": matched}); err != nil {
		return err
	}
	if err := w.subs.InsertSubmissionRows(ctx, sub.ID, classified.Rows, classified.Values); err != nil {
		return err
	}
	return w.jobs.CompleteJob(ctx, job.ID, progressJSON("selesai", 100, map[string]any{
		"rows": len(classified.Rows), "matched": matched, "need_review": len(classified.Rows) - matched,
	}), time.Now().UTC().Format(time.RFC3339))
}

func (w *Worker) progress(ctx context.Context, jobID, step string, percent int, extra map[string]any) error {
	return w.jobs.UpdateJobProgress(ctx, jobID, progressJSON(step, percent, extra))
}

func progressJSON(step string, percent int, extra map[string]any) string {
	m := map[string]any{"step": step, "percent": percent}
	for k, v := range extra {
		m[k] = v
	}
	b, _ := json.Marshal(m)
	return string(b)
}
