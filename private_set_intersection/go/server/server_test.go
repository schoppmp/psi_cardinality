package server

import (
	"fmt"
	// "github.com/openmined/psi/client"
	psi_ds "github.com/openmined/psi/datastructure"
	psi_proto "github.com/openmined/psi/pb"
	// "google.golang.org/protobuf/proto"
	"testing"
)

var paramsTable = []struct {
	revealIntersection bool
	ds                 psi_ds.DataStructure
}{
	{revealIntersection: true, ds: psi_ds.Raw},
	{revealIntersection: true, ds: psi_ds.Gcs},
	{revealIntersection: true, ds: psi_ds.BloomFilter},
	{revealIntersection: false, ds: psi_ds.Raw},
	{revealIntersection: false, ds: psi_ds.Gcs},
	{revealIntersection: false, ds: psi_ds.BloomFilter},
}

var inputsTable = []struct {
	fpr    float64
	inputs int
}{
	{fpr: 0.01, inputs: 1},
	{fpr: 0.01, inputs: 10},
	{fpr: 0.01, inputs: 100},
	{fpr: 0.01, inputs: 1000},
	{fpr: 0.000001, inputs: 1},
	{fpr: 0.000001, inputs: 10},
	{fpr: 0.000001, inputs: 100},
	{fpr: 0.000001, inputs: 1000},
}

func BenchmarkServerSetup(b *testing.B) {
	for _, pt := range paramsTable {
		for _, it := range inputsTable {
			b.Run(fmt.Sprintf("reveal_%v_ds_%s_fpr_%f_inputs_%d", pt.revealIntersection, pt.ds, it.fpr, it.inputs), func(b *testing.B) {
				b.ReportAllocs()
				total := 0

				server, err := CreateWithNewKey(pt.revealIntersection)
				if err != nil || server == nil {
					b.Errorf("failed to get server")
				}

				serverInputs := []string{}
				for i := 0; i < it.inputs; i++ {
					serverInputs = append(serverInputs, "Element "+string(i))
				}

				for i := 0; i < b.N; i++ {
					_, err := server.CreateSetupMessage(it.fpr, int64(10000), serverInputs, pt.ds)
					if err != nil {
						b.Errorf("failed to create setup msg %v", err)
					}
					total += it.inputs
				}

				b.ReportMetric(float64(total), "ElementsProcessed")
			})
		}
	}
}

// const fpr3 = 0.001
// const fpr6 = 0.000001

// func BenchmarkServerSetup1fpr3(b *testing.B)               { benchmarkServerSetup(1, fpr3, false, b) }
// func BenchmarkServerSetup10fpr3(b *testing.B)              { benchmarkServerSetup(10, fpr3, false, b) }
// func BenchmarkServerSetup100fpr3(b *testing.B)             { benchmarkServerSetup(100, fpr3, false, b) }
// func BenchmarkServerSetup1000fpr3(b *testing.B)            { benchmarkServerSetup(1000, fpr3, false, b) }
// func BenchmarkServerSetup10000fpr3(b *testing.B)           { benchmarkServerSetup(10000, fpr3, false, b) }
// func BenchmarkServerSetupIntersection1fpr3(b *testing.B)   { benchmarkServerSetup(1, fpr3, true, b) }
// func BenchmarkServerSetupIntersection10fpr3(b *testing.B)  { benchmarkServerSetup(10, fpr3, true, b) }
// func BenchmarkServerSetupIntersection100fpr3(b *testing.B) { benchmarkServerSetup(100, fpr3, true, b) }
// func BenchmarkServerSetupIntersection1000fpr3(b *testing.B) {
// 	benchmarkServerSetup(1000, fpr3, true, b)
// }
// func BenchmarkServerSetupIntersection10000fpr3(b *testing.B) {
// 	benchmarkServerSetup(10000, fpr3, true, b)
// }

// func BenchmarkServerSetup1fpr6(b *testing.B)               { benchmarkServerSetup(1, fpr6, false, b) }
// func BenchmarkServerSetup10fpr6(b *testing.B)              { benchmarkServerSetup(10, fpr6, false, b) }
// func BenchmarkServerSetup100fpr6(b *testing.B)             { benchmarkServerSetup(100, fpr6, false, b) }
// func BenchmarkServerSetup1000fpr6(b *testing.B)            { benchmarkServerSetup(1000, fpr6, false, b) }
// func BenchmarkServerSetup10000fpr6(b *testing.B)           { benchmarkServerSetup(10000, fpr6, false, b) }
// func BenchmarkServerSetupIntersection1fpr6(b *testing.B)   { benchmarkServerSetup(1, fpr6, true, b) }
// func BenchmarkServerSetupIntersection10fpr6(b *testing.B)  { benchmarkServerSetup(10, fpr6, true, b) }
// func BenchmarkServerSetupIntersection100fpr6(b *testing.B) { benchmarkServerSetup(100, fpr6, true, b) }
// func BenchmarkServerSetupIntersection1000fpr6(b *testing.B) {
// 	benchmarkServerSetup(1000, fpr6, true, b)
// }
// func BenchmarkServerSetupIntersection10000fpr6(b *testing.B) {
// 	benchmarkServerSetup(10000, fpr6, true, b)
// }

// var dummyResponse *psi_proto.Response

// func benchmarkServerProcessRequest(cnt int, revealIntersection bool, b *testing.B) {
// 	b.ReportAllocs()
// 	total := 0

// 	client, err := client.CreateWithNewKey(revealIntersection)
// 	if err != nil || client == nil {
// 		b.Errorf("failed to get client")
// 	}
// 	server, err := CreateWithNewKey(revealIntersection)
// 	if err != nil || server == nil {
// 		b.Errorf("failed to get server")
// 	}

// 	inputs := []string{}
// 	for i := 0; i < cnt; i++ {
// 		inputs = append(inputs, "Element "+string(i))
// 	}

// 	request, err := client.CreateRequest(inputs)
// 	if err != nil {
// 		b.Errorf("failed to create request %v", err)
// 	}

// 	for n := 0; n < b.N; n++ {
// 		serverResp, err := server.ProcessRequest(request)
// 		if err != nil {
// 			b.Errorf("failed to process request %v", err)
// 		}
// 		total += cnt
// 		b.ReportMetric(float64(proto.Size(serverResp)), "ResponseSize")
// 		//ugly hack for preventing compiler optimizations
// 		dummyResponse = serverResp
// 	}
// 	b.ReportMetric(float64(total), "ElementsProcessed")
// }

// func BenchmarkServerProcessRequest1(b *testing.B)     { benchmarkServerProcessRequest(1, false, b) }
// func BenchmarkServerProcessRequest10(b *testing.B)    { benchmarkServerProcessRequest(10, false, b) }
// func BenchmarkServerProcessRequest100(b *testing.B)   { benchmarkServerProcessRequest(100, false, b) }
// func BenchmarkServerProcessRequest1000(b *testing.B)  { benchmarkServerProcessRequest(1000, false, b) }
// func BenchmarkServerProcessRequest10000(b *testing.B) { benchmarkServerProcessRequest(10000, false, b) }
// func BenchmarkServerProcessRequestIntersection1(b *testing.B) {
// 	benchmarkServerProcessRequest(1, true, b)
// }
// func BenchmarkServerProcessRequestIntersection10(b *testing.B) {
// 	benchmarkServerProcessRequest(10, true, b)
// }
// func BenchmarkServerProcessRequestIntersection100(b *testing.B) {
// 	benchmarkServerProcessRequest(100, true, b)
// }
// func BenchmarkServerProcessRequestIntersection1000(b *testing.B) {
// 	benchmarkServerProcessRequest(1000, true, b)
// }
// func BenchmarkServerProcessRequestIntersection10000(b *testing.B) {
// 	benchmarkServerProcessRequest(10000, true, b)
// }
