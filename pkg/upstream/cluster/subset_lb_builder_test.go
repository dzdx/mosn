package cluster

import (
	"github.com/stretchr/testify/require"
	"mosn.io/mosn/pkg/types"
	"testing"
)

func TestSubsetLbBuilder_InitIndex(t *testing.T) {
	subsetInfo := NewLBSubsetInfo(exampleSubsetConfig())
	hosts := exampleHostConfigs()
	builder := NewSubsetLbBuilder(&clusterInfo{}, createHostset(hosts), subsetMergeKeys(subsetInfo.SubsetKeys(), subsetInfo.DefaultSubset()))
	merged := make(map[string]map[string][]int)
	for i, host := range hosts {
		for _, key := range subsetMergeKeys(subsetInfo.SubsetKeys(), subsetInfo.DefaultSubset()) {
			val, ok := host.MetaData[key]
			if !ok {
				continue
			}
			if _, ok := merged[key]; !ok {
				merged[key] = make(map[string][]int)
			}
			merged[key][val] = append(merged[key][val], i)
		}
	}
	require.Len(t, builder.indexer, len(merged))
	for key, vals := range merged {
		require.Len(t, builder.indexer[key], len(vals))
		for val, hosts := range vals {
			indexArray := builder.indexer[key][val].AppendTo(nil)
			require.Equal(t, hosts, indexArray)
		}
	}
}

func TestSubsetLbBuilder_ExtractKvComb(t *testing.T) {
	subsetInfo := NewLBSubsetInfo(exampleSubsetConfig())
	builder := NewSubsetLbBuilder(&clusterInfo{}, createHostset(exampleHostConfigs()), subsetMergeKeys(subsetInfo.SubsetKeys(), subsetInfo.DefaultSubset()))
	var kvList []types.SubsetMetadata
	kvList = builder.extractKvComb([]string{"version"})
	// [[{version 1.0}] [{version 1.1}] [{version 1.2-pre}]]
	require.Len(t, kvList, 3)
	for _, kvs := range kvList {
		require.Len(t, kvs, 1)
	}
	kvList = builder.extractKvComb([]string{"stage", "type"})
	// [[{stage prod} {type bigmem}] [{stage prod} {type std}] [{stage dev} {type std}] [{stage dev} {type bigmem}]]
	require.Len(t, kvList, 4)
	for _, kvs := range kvList {
		require.Len(t, kvs, 2)
	}

	kvList = builder.extractKvComb([]string{"stage", "version"})
	// [[{stage prod} {version 1.0}] [{stage prod} {version 1.1}] [{stage prod} {version 1.2-pre}] [{stage dev} {version 1.0}] [{stage dev} {version 1.1}] [{stage dev} {version 1.2-pre}]]
	require.Len(t, kvList, 6)
	for _, kvs := range kvList {
		require.Len(t, kvs, 2)
	}
}

func TestSubsetMergeKeys(t *testing.T) {
	keys := subsetMergeKeys([]types.SortedStringSetType{
		types.InitSet([]string{"1", "2", "3"}),
		types.InitSet([]string{"2", "3", "4"}),
	}, []types.Pair{
		{
			T1: "5",
			T2: "6",
		},
	})
	require.Len(t, keys, 5)
}
