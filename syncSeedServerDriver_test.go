package main

import (
	"fmt"
	"seed-sync/driver/client"
	seedSyncModel "seed-sync/model/seed-sync"
	"testing"
)

func TestSyncSeed(t *testing.T) {
	var request = &seedSyncModel.SeedSyncRequest{
		Sites: []string{"hdhome", "hhan", "cyanbug", "icc2022", "ptvicomo", "soulvoice", "btschool", "qingwapt", "hdfans", "carpt", "hdkyl", "ubits", "audiences", "1ptba", "crabpt"},
		Torrents: []seedSyncModel.TorrentForSeedSyncRequest{
			{
				PiecesHash: "f1fa809ed5706d44e38cf02a54022f0a785dcae2",
				FilesHash:  "9793bb9da92eb74a7bb4214980d1c6ff09e58e97",
			},
		},
	}
	result, err := client.SeedSyncServerClient.SyncSeed(request)
	if err != nil {
		t.Errorf("SyncSeed error: %v", err)
	}
	fmt.Println(result)
}
