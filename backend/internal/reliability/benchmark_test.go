package reliability_test

import (
	"context"
	"testing"

	backupModel "github.com/cloudstorex/backend/internal/backup/model"
	backupSvc "github.com/cloudstorex/backend/internal/backup/service"
	drModel "github.com/cloudstorex/backend/internal/disasterrecovery/model"
	drSvc "github.com/cloudstorex/backend/internal/disasterrecovery/service"
	replicationModel "github.com/cloudstorex/backend/internal/replication/model"
	replicationSvc "github.com/cloudstorex/backend/internal/replication/service"
	repairModel "github.com/cloudstorex/backend/internal/selfhealing/model"
	repairSvc "github.com/cloudstorex/backend/internal/selfhealing/service"
	verifModel "github.com/cloudstorex/backend/internal/verification/model"
	verifSvc "github.com/cloudstorex/backend/internal/verification/service"

	"github.com/google/uuid"
)

// Mock Repositories for Benchmarks
type mockReplicationRepo struct{}
func (m *mockReplicationRepo) Create(ctx context.Context, repl *replicationModel.ObjectReplication) error { return nil }
func (m *mockReplicationRepo) Update(ctx context.Context, repl *replicationModel.ObjectReplication) error { return nil }
func (m *mockReplicationRepo) GetByID(ctx context.Context, id uuid.UUID) (*replicationModel.ObjectReplication, error) { return nil, nil }
func (m *mockReplicationRepo) ListByObjectID(ctx context.Context, objectID uuid.UUID) ([]replicationModel.ObjectReplication, error) { return nil, nil }
func (m *mockReplicationRepo) ListByStatus(ctx context.Context, status replicationModel.ReplicationStatus, limit int) ([]replicationModel.ObjectReplication, error) { return nil, nil }

type mockBackupRepo struct{}
func (m *mockBackupRepo) CreateJob(ctx context.Context, job *backupModel.BackupJob) error { return nil }
func (m *mockBackupRepo) UpdateJob(ctx context.Context, job *backupModel.BackupJob) error { return nil }
func (m *mockBackupRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*backupModel.BackupJob, error) { return &backupModel.BackupJob{}, nil }
func (m *mockBackupRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]backupModel.BackupJob, error) { return nil, nil }
func (m *mockBackupRepo) CreateRecord(ctx context.Context, rec *backupModel.BackupRecord) error { return nil }
func (m *mockBackupRepo) ListRecordsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]backupModel.BackupRecord, error) { return nil, nil }

type mockDRRepo struct{}
func (m *mockDRRepo) CreatePlan(ctx context.Context, plan *drModel.RecoveryPlan) error { return nil }
func (m *mockDRRepo) GetPlanByID(ctx context.Context, id uuid.UUID) (*drModel.RecoveryPlan, error) { return &drModel.RecoveryPlan{}, nil }
func (m *mockDRRepo) ListPlansByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]drModel.RecoveryPlan, error) { return nil, nil }
func (m *mockDRRepo) CreateJob(ctx context.Context, job *drModel.RecoveryJob) error { return nil }
func (m *mockDRRepo) UpdateJob(ctx context.Context, job *drModel.RecoveryJob) error { return nil }
func (m *mockDRRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*drModel.RecoveryJob, error) { return &drModel.RecoveryJob{}, nil }
func (m *mockDRRepo) ListJobsByPlan(ctx context.Context, planID uuid.UUID) ([]drModel.RecoveryJob, error) { return nil, nil }

type mockVerifRepo struct{}
func (m *mockVerifRepo) CreateJob(ctx context.Context, job *verifModel.ConsistencyJob) error { return nil }
func (m *mockVerifRepo) UpdateJob(ctx context.Context, job *verifModel.ConsistencyJob) error { return nil }
func (m *mockVerifRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*verifModel.ConsistencyJob, error) { return &verifModel.ConsistencyJob{}, nil }
func (m *mockVerifRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]verifModel.ConsistencyJob, error) { return nil, nil }
func (m *mockVerifRepo) CreateDiscrepancy(ctx context.Context, rec *verifModel.DiscrepancyRecord) error { return nil }
func (m *mockVerifRepo) ListDiscrepanciesByWorkspace(ctx context.Context, workspaceID uuid.UUID, unresolvedOnly bool) ([]verifModel.DiscrepancyRecord, error) { return nil, nil }
func (m *mockVerifRepo) MarkDiscrepancyResolved(ctx context.Context, id uuid.UUID) error { return nil }

type mockRepairRepo struct{}
func (m *mockRepairRepo) CreateJob(ctx context.Context, job *repairModel.RepairJob) error { return nil }
func (m *mockRepairRepo) UpdateJob(ctx context.Context, job *repairModel.RepairJob) error { return nil }
func (m *mockRepairRepo) GetJobByID(ctx context.Context, id uuid.UUID) (*repairModel.RepairJob, error) { return &repairModel.RepairJob{}, nil }
func (m *mockRepairRepo) ListJobsByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]repairModel.RepairJob, error) { return nil, nil }

func BenchmarkReplicationScheduling(b *testing.B) {
	svc := replicationSvc.NewService(&mockReplicationRepo{}, nil, nil, nil)
	ctx := context.Background()
	req := replicationSvc.ScheduleRequest{
		WorkspaceID:    uuid.New(),
		BucketID:       uuid.New(),
		ObjectID:       uuid.New(),
		ObjectKey:      "data/file.bin",
		PrimaryProvider: "aws-s3-east",
		ReplicaProvider: "aws-s3-west",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ScheduleReplication(ctx, req)
	}
}

func BenchmarkBackupScheduling(b *testing.B) {
	svc := backupSvc.NewService(&mockBackupRepo{}, nil, nil, nil)
	ctx := context.Background()
	req := backupSvc.ScheduleBackupRequest{
		WorkspaceID:    uuid.New(),
		BucketID:       uuid.New(),
		SourceProvider: "aws-s3-east",
		TargetProvider: "aws-glacier",
		BackupType:     backupModel.BackupTypeFull,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ScheduleBackup(ctx, req)
	}
}

func BenchmarkRestoreScheduling(b *testing.B) {
	svc := drSvc.NewService(&mockDRRepo{}, nil, nil, nil)
	ctx := context.Background()
	req := drSvc.TriggerRecoveryRequest{
		PlanID:      uuid.New(),
		PointInTime: nil,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.TriggerRecovery(ctx, req)
	}
}

func BenchmarkVerificationScheduling(b *testing.B) {
	svc := verifSvc.NewService(&mockVerifRepo{}, nil, nil, nil)
	ctx := context.Background()
	req := verifSvc.ScheduleCheckRequest{
		WorkspaceID: uuid.New(),
		BucketID:    uuid.New(),
		ProviderID:  "aws-s3-east",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ScheduleVerification(ctx, req)
	}
}

func BenchmarkSelfHealingScheduling(b *testing.B) {
	svc := repairSvc.NewService(&mockRepairRepo{}, nil, nil, nil)
	ctx := context.Background()
	req := repairSvc.ScheduleRepairRequest{
		WorkspaceID:    uuid.New(),
		BucketID:       uuid.New(),
		ObjectID:       uuid.New(),
		ObjectKey:      "data/file.bin",
		TargetProvider: "aws-s3-east",
		SourceProvider: "aws-s3-west",
		SourceType:     repairModel.SourceReplica,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.ScheduleRepair(ctx, req)
	}
}
