package disk

import (
	pollster "metricly/internal/collector"
	helper "metricly/internal/pollster/tests"
	"os"
	"testing"
)

func TestGetMountPoints(t *testing.T) {

	mntContent := `/dev/mapper/luks-49c47969-6ea3-4aaa-8200-9768d072c21c / btrfs rw,seclabel,relatime,compress=zstd:1,ssd,discard=async,space_cache=v2,subvolid=257,subvol=/root 0 0
devtmpfs /dev devtmpfs rw,seclabel,nosuid,size=4096k,nr_inodes=4063587,mode=755,inode64 0 0
tmpfs /dev/shm tmpfs rw,seclabel,nosuid,nodev,inode64 0 0
devpts /dev/pts devpts rw,seclabel,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=000 0 0
sysfs /sys sysfs rw,seclabel,nosuid,nodev,noexec,relatime 0 0
/dev/nvme0n1p2 /boot ext4 rw,seclabel,relatime 0 0`

	collectorSource := "mounts_test.txt"
	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procMounts = collectorSource

	// start testing target function
	mounts, err := getMountPoints()
	if err != nil {
		t.Fatal(err)
	}
	if len(mounts) != 5 {
		t.Errorf("incorrect count of mounts: %v", mounts)
	}

	expectedMounts := []string{"/", "/dev", "/dev/pts", "/sys", "/boot"}
	for i := range len(mounts) {
		if mounts[i] != expectedMounts[i] {
			t.Error("incorrect mount")
		}
	}
}

func TestParseDiskStats(t *testing.T) {
	// Mock /proc/diskstats content
	collectorSource := "diskstats.txt"
	mntContent := `8       0 sda 157698 987 4056738 364879 45893 123 987235 456812 0 45601 45601
	   8       1 sda1 10045 64 405678 100 4568 0 12345 45678 0 123 123
	   8       16 sdb 250698 587 2056738 264879 25893 53 287235 256812 0 25601 25601`

	err := helper.SetupCollectorSources(collectorSource, mntContent)
	if err != nil {
		t.Fatalf("failed to setup collector file: %v", err)
	}
	defer os.Remove(collectorSource)
	procDiskStats = collectorSource

	// start testing target function
	mapDiskStats, err := parseDiskStats()
	if err != nil {
		t.Fatalf("failed to parse disk stats: %s", err)
	}

	if mapDiskStats["sda"].ReadsCompleted != 157698 {
		t.Errorf("expected ReadsCompleted=157698, got %d", mapDiskStats["sda"].ReadsCompleted)
	}
	if mapDiskStats["sda1"].IOInProgress != 0 {
		t.Errorf("expected IOInProgress=0, got %d", mapDiskStats["sdb"].IOInProgress)
	}
	if mapDiskStats["sdb"].ReadThroughputBytes != 1053049856 {
		t.Errorf("expected ReadsCompleted=1053049856, got %d", mapDiskStats["sda"].ReadThroughputBytes)
	}

	mc := pollster.CreateMetricCollector()
	RegisterDiskMetrics(mc)
	ReportDiskUsage(mc)

	if metric, ok := mc.Data["disk_reads_completed_total|sda"]; !ok || metric.Value != 157698 {
		t.Errorf("unexpected value for disk_reads_completed_total: got %v, want %v", metric.Value, 157698)
	}
	if metric, ok := mc.Data["disk_io_in_progress|sda1"]; !ok || metric.Value != 0 {
		t.Errorf("unexpected value for disk_reads_completed_total: got %v, want %v", metric.Value, 0)
	}
	if metric, ok := mc.Data["disk_read_throughput_bytes|sdb"]; !ok || metric.Value != 1053049856 {
		t.Errorf("unexpected value for disk_reads_completed_total: got %v, want %v", metric.Value, 1053049856)
	}

}
