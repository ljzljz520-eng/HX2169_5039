package domain

import "testing"

func TestRecordLifecycle(t *testing.T) {
	record := NewRecord("r1", "开场舞", "舞蹈队", "舞蹈", "候选", 1)
	var err error
	record, err = BeginReview(record, 2)
	if err != nil || record.Status != StatusReview {
		t.Fatalf("begin review: %#v %v", record, err)
	}
	record, err = ReviewRecord(record, 90, 3)
	if err != nil || record.Status != StatusApproved {
		t.Fatalf("review: %#v %v", record, err)
	}
	record, err = ArchiveRecord(record, 4)
	if err != nil || record.Status != StatusArchived {
		t.Fatalf("archive: %#v %v", record, err)
	}
}

func TestScoreBands(t *testing.T) {
	if BandForScore(82, DefaultBands()).Name != "silver" {
		t.Fatal("unexpected band")
	}
	if WeightedScore(80, 70, 60) != 73 {
		t.Fatal("unexpected weighted score")
	}
}
