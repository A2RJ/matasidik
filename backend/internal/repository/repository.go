package repository

import (
	"context"

	"matasidik/backend/internal/domain"
)

type UserRepository interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.UserAuth, error)
	ListPositions(ctx context.Context) ([]domain.Position, error)
	CreatePosition(ctx context.Context, p *domain.Position) error
	UpdatePosition(ctx context.Context, p *domain.Position) error
	CreateUser(ctx context.Context, u *domain.User, passwordHash string) error
	UpdateUser(ctx context.Context, u *domain.User) error
	UpdateUserName(ctx context.Context, id, nama string) error
	UpdateUserPassword(ctx context.Context, id, passwordHash string) error
	CreateFaskes(ctx context.Context, f *domain.Faskes) error
	UpdateFaskes(ctx context.Context, f *domain.Faskes) error
}

type AuthRepository interface {
	CreateSession(ctx context.Context, s *domain.AuthSession) error
	GetSessionByHash(ctx context.Context, hash string) (*domain.AuthSession, error)
	RevokeSession(ctx context.Context, id, revokedAt string) error
	RevokeSessionsForUser(ctx context.Context, userID, revokedAt string) error
	UpdateLastLogin(ctx context.Context, userID, at string) error
}

type MasterRepository interface {
	ListWilayah(ctx context.Context) ([]domain.Wilayah, error)
	ListFaskes(ctx context.Context) ([]domain.Faskes, error)
	GetFaskes(ctx context.Context, id string) (*domain.Faskes, error)
	ListBidang(ctx context.Context) ([]domain.Bidang, error)
	ListPeriods(ctx context.Context) ([]domain.Periode, error)
	GetPeriodByKode(ctx context.Context, kode string) (*domain.Periode, error)
	GetPeriodByID(ctx context.Context, id string) (*domain.Periode, error)
	DefaultPeriodID(ctx context.Context) (string, error)
	GetPublishedRelease(ctx context.Context, tahun int) (*domain.MasterRelease, error)
	GetMasterRelease(ctx context.Context, id string) (*domain.MasterRelease, error)
	GetWorkflowVersion(ctx context.Context, id string) (*domain.WorkflowVersion, error)

	ListLabels(ctx context.Context) ([]domain.Label, error)
	GetLabelByID(ctx context.Context, id string) (*domain.Label, error)
	GetLabelByKode(ctx context.Context, kode, releaseID string) (*domain.Label, error)
	ListIndicators(ctx context.Context) ([]domain.Indikator, error)
	GetIndicatorByID(ctx context.Context, id string) (*domain.Indikator, error)
	GetIndicatorByKode(ctx context.Context, kode string) (*domain.Indikator, error)
	ListIndicatorLabels(ctx context.Context, indikatorID string) ([]domain.IndikatorLabel, error)
	ListIndicatorLabelsForLabel(ctx context.Context, labelID string) ([]domain.IndikatorLabel, error)
	ListIKUComponents(ctx context.Context, parentLabelID string) ([]domain.IKUComponent, error)
	ListTargets(ctx context.Context, periodeID string) ([]domain.Target, error)
	ListCapaianTahunan(ctx context.Context, tahun int, scopeType, scopeID string) ([]domain.CapaianTahunan, error)
	GetCapaianTahunan(ctx context.Context, id string) (*domain.CapaianTahunan, error)
	UpsertCapaianTahunan(ctx context.Context, c *domain.CapaianTahunan) error
	DeleteCapaianTahunan(ctx context.Context, id string) error
	ListDashboardRows(ctx context.Context) ([]domain.DashboardRow, error)
	GetDashboardRowByKode(ctx context.Context, kode string) (*domain.DashboardRow, error)

	CreateIndicator(ctx context.Context, indikator *domain.Indikator, labelIDs []string) error
	CreateLabel(ctx context.Context, label *domain.Label) error
	UpsertTarget(ctx context.Context, target *domain.Target) error
}

type SubmissionRepository interface {
	FindIdempotentSubmission(ctx context.Context, key string) (*domain.Submission, error)
	CreateUpload(ctx context.Context, p domain.CreateUploadParams) (*domain.Submission, *domain.Job, error)
	GetSubmission(ctx context.Context, id string) (*domain.Submission, error)
	GetSubmissionByKode(ctx context.Context, kode string) (*domain.Submission, error)
	GetFileObject(ctx context.Context, id string) (*domain.FileObject, error)
	GetSubmissionByFileID(ctx context.Context, fileID string) (*domain.Submission, error)
	ListSubmissions(ctx context.Context, f domain.SubmissionFilter) ([]domain.Submission, error)
	GetRows(ctx context.Context, submissionID string) ([]domain.SubmissionRow, error)
	GetRow(ctx context.Context, rowID string) (*domain.SubmissionRow, error)
	UpdateRowClassification(ctx context.Context, rowID, labelID string, status domain.ClassificationStatus, errMsg string) error
	UpsertIndicatorValueForRow(ctx context.Context, row *domain.SubmissionRow, sub *domain.Submission) error
	UpdateSubmissionStatus(ctx context.Context, id string, status domain.SubmissionStatus, step domain.WorkflowStep) error
	CreateReview(ctx context.Context, review *domain.SubmissionReview) error
	ListReviews(ctx context.Context, submissionID string) ([]domain.SubmissionReview, error)
	ListRevisions(ctx context.Context, submissionID string) ([]domain.SubmissionRevision, error)
	InsertSubmissionRows(ctx context.Context, submissionID string, rows []domain.SubmissionRowInput, values []domain.IndicatorValueInput) error
	FindDuplicateRow(ctx context.Context, submissionID, labelID, indikatorID, scopeCode string) (bool, error)
	CountBlockingRows(ctx context.Context, submissionID string) (int, error)
}

type JobRepository interface {
	GetJob(ctx context.Context, id string) (*domain.Job, error)
	ClaimJob(ctx context.Context, workerID, now string) (*domain.Job, error)
	UpdateJobProgress(ctx context.Context, id, progressJSON string) error
	CompleteJob(ctx context.Context, id, progressJSON, finishedAt string) error
	FailJob(ctx context.Context, id, errMsg, finishedAt string) error
	RequeueStaleJobs(ctx context.Context, before string) error
}

type ValueRepository interface {
	ListFinalIndicatorValues(ctx context.Context, periodeID string) ([]domain.IndicatorValue, error)
	ListFinalIndicatorValuesForScope(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IndicatorValue, error)
	GetIKUValue(ctx context.Context, labelID, periodeID, scopeType, scopeID string) (*domain.IKUValue, error)
	ListIKUValues(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error)
	ListIKUValueDetails(ctx context.Context, ikuValueID string) ([]domain.IKUValueDetail, error)
	UpsertIKUValues(ctx context.Context, writes []domain.IKUValueWrite) error
	CreateCalculationRun(ctx context.Context, run *domain.CalculationRun) error
	FinishCalculationRun(ctx context.Context, id, status, errMsg, finishedAt string) error
	LatestCalculationRun(ctx context.Context, periodeID string) (*domain.CalculationRun, error)
	ListFinalFileHashes(ctx context.Context, periodeID, scopeType, scopeID string) ([]string, error)
	GetIndicatorValueBySubmissionRow(ctx context.Context, submissionRowID string) (*domain.IndicatorValue, error)
	DashboardSPM(ctx context.Context, f domain.DashboardFilter) ([]domain.DashboardSPMItem, error)
	TraceSPM(ctx context.Context, kode string, f domain.DashboardFilter) (*domain.SPMTrace, error)
}

type FactorRepository interface {
	ListIKUValuesForFactors(ctx context.Context, periodeID, scopeType, scopeID string) ([]domain.IKUValue, error)
	GetFaktor(ctx context.Context, ikuValueID string) (*domain.Faktor, error)
	UpsertFaktor(ctx context.Context, f *domain.Faktor) error
	ListTindakLanjut(ctx context.Context, periodeID string) ([]domain.TindakLanjut, error)
	CreateTindakLanjut(ctx context.Context, tl *domain.TindakLanjut) error
	UpdateTindakLanjut(ctx context.Context, tl *domain.TindakLanjut) error
	GetTindakLanjut(ctx context.Context, id string) (*domain.TindakLanjut, error)
}

type AuditRepository interface {
	InsertAudit(ctx context.Context, log *domain.AuditLog) error
	ListAudit(ctx context.Context, f domain.AuditFilter) ([]domain.AuditLog, error)
}

type Repository interface {
	UserRepository
	AuthRepository
	MasterRepository
	SubmissionRepository
	JobRepository
	ValueRepository
	FactorRepository
	AuditRepository
}
