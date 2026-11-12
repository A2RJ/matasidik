package domain

// Role adalah peran pengguna pada pilot MATA SIDIK.
// Auth produksi belum diaktifkan; peran tetap dipakai untuk otorisasi aksi.
type Role string

const (
	RolePKMOperator Role = "pkm_operator"
	RolePKMPJ       Role = "pkm_pj"
	RoleVerifier    Role = "dinas_verifier"
	RoleKabid       Role = "kabid"
	RoleSekretaris  Role = "sekretaris"
	RoleIKUOwner    Role = "iku_owner"
	RolePimpinan    Role = "pimpinan"
	RoleAdmin       Role = "admin"
)

var AllRoles = []Role{
	RolePKMOperator, RolePKMPJ, RoleVerifier, RoleKabid,
	RoleSekretaris, RoleIKUOwner, RolePimpinan, RoleAdmin,
}

func (r Role) IsApprover() bool {
	return r == RoleVerifier || r == RoleKabid || r == RoleSekretaris || r == RoleAdmin
}

// SubmissionStatus mengikuti tujuh status formal di docs/rules.md §6.1.
type SubmissionStatus string

const (
	StatusDraf                    SubmissionStatus = "draf"
	StatusDikirim                 SubmissionStatus = "dikirim"
	StatusDiverifikasi            SubmissionStatus = "diverifikasi"
	StatusDikembalikanPetugas     SubmissionStatus = "dikembalikan_petugas"
	StatusDikembalikanVerifikator SubmissionStatus = "dikembalikan_verifikator"
	StatusFinal                   SubmissionStatus = "final"
	StatusDiubahKhusus            SubmissionStatus = "diubah_khusus"
)

// WorkflowStep adalah nama langkah internal yang dipin ke submission.
type WorkflowStep string

const (
	StepDraftPKM      WorkflowStep = "DRAFT_PKM"
	StepVerifyDinas   WorkflowStep = "VERIFY_DINAS"
	StepValidateKabid WorkflowStep = "VALIDATE_KABID"
	StepFinal         WorkflowStep = "FINAL"
	StepSpecialChange WorkflowStep = "SPECIAL_CHANGE"
)

// ClassificationStatus adalah hasil klasifikasi deterministik.
type ClassificationStatus string

const (
	ClassMatched      ClassificationStatus = "MATCHED"
	ClassMismatch     ClassificationStatus = "MISMATCH"
	ClassUnknownIKU   ClassificationStatus = "UNKNOWN_IKU"
	ClassUnknownInd   ClassificationStatus = "UNKNOWN_INDICATOR"
	ClassRollupOnly   ClassificationStatus = "IKU_ROLLUP_ONLY"
	ClassAmbiguous    ClassificationStatus = "AMBIGUOUS"
	ClassDuplicate    ClassificationStatus = "DUPLICATE"
	ClassInvalidValue ClassificationStatus = "INVALID_VALUE"
	ClassScopeDenied  ClassificationStatus = "SCOPE_DENIED"
)

// IsBlocking menandai baris yang menghalangi submit.
func (c ClassificationStatus) IsBlocking() bool {
	switch c {
	case ClassUnknownIKU, ClassUnknownInd, ClassRollupOnly, ClassMismatch,
		ClassAmbiguous, ClassDuplicate, ClassInvalidValue, ClassScopeDenied:
		return true
	default:
		return false
	}
}

// ValueStatus adalah status kelengkapan nilai IKU.
type ValueStatus string

const (
	ValueComplete ValueStatus = "COMPLETE"
	ValuePartial  ValueStatus = "PARTIAL"
	ValueNoData   ValueStatus = "NO_DATA"
	ValueStale    ValueStatus = "STALE"
)

// AggregationMethod adalah metode agregasi IKU.
type AggregationMethod string

const (
	AggRatioOfSums     AggregationMethod = "RATIO_OF_SUMS"
	AggSum             AggregationMethod = "SUM"
	AggWeightedAverage AggregationMethod = "WEIGHTED_AVERAGE"
	AggFormula         AggregationMethod = "FORMULA"
	AggNone            AggregationMethod = "NONE"
)

// TLStatus adalah status tindak lanjut.
type TLStatus string

const (
	TLOpen       TLStatus = "open"
	TLInProgress TLStatus = "in_progress"
	TLDone       TLStatus = "done"
	TLCancelled  TLStatus = "cancelled"
)

// JobStatus adalah status antrean job SQLite.
type JobStatus string

const (
	JobQueued  JobStatus = "queued"
	JobRunning JobStatus = "running"
	JobDone    JobStatus = "done"
	JobFailed  JobStatus = "failed"
	JobDead    JobStatus = "dead"
)

var AllSubmissionStatuses = []SubmissionStatus{
	StatusDraf, StatusDikirim, StatusDiverifikasi,
	StatusDikembalikanPetugas, StatusDikembalikanVerifikator,
	StatusFinal, StatusDiubahKhusus,
}

var AllClassificationStatuses = []ClassificationStatus{
	ClassMatched, ClassMismatch, ClassUnknownIKU, ClassUnknownInd,
	ClassRollupOnly, ClassAmbiguous, ClassDuplicate, ClassInvalidValue, ClassScopeDenied,
}

var AllValueStatuses = []ValueStatus{ValueComplete, ValuePartial, ValueNoData, ValueStale}

var AllTLStatuses = []TLStatus{TLOpen, TLInProgress, TLDone, TLCancelled}
