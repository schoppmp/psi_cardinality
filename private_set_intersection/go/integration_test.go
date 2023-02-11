package integration_test

import (
	psi_client "github.com/openmined/psi/client"
	psi_ds "github.com/openmined/psi/datastructure"
	psi_proto "github.com/openmined/psi/pb"
	psi_server "github.com/openmined/psi/server"
	psi_version "github.com/openmined/psi/version"
	"google.golang.org/protobuf/proto"
	"regexp"
	"testing"
)

var clientKey = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
var serverKey = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

const fpr = 0.001
const numClientInputs = 10
const numServerInputs = 100

func TestVersion(t *testing.T) {
	version := psi_version.Version()
	r, _ := regexp.Compile(`[0-9]+[.][0-9]+[.][0-9]+(-[A-Za-z0-9]+)?`)

	matched := r.MatchString(version)
	if !matched {
		t.Errorf("Invalid version: %v", version)
	}

	client, err := psi_client.CreateWithNewKey(false)
	if err != nil {
		t.Errorf("Failed to create a PSI client %v", err)
	}
	server, err := psi_server.CreateWithNewKey(false)
	if err != nil {
		t.Errorf("Failed to create a PSI server %v", err)
	}

	clientVersion := client.Version()
	serverVersion := server.Version()
	matched = r.MatchString(clientVersion)
	if !matched {
		t.Errorf("Invalid client version: %v", version)
	}
	matched = r.MatchString(serverVersion)
	if !matched {
		t.Errorf("Invalid server version: %v", version)
	}

	if version != clientVersion || version != serverVersion || clientVersion != serverVersion {
		t.Errorf("Missmatched versions pkg (%s), client (%s), server (%s)", version, clientVersion, serverVersion)
	}

}
func TestIntegrationIntersection(t *testing.T) {
	testCases := []struct {
		revealIntersection bool
		ds                 psi_ds.DataStructure
	}{
		{true, psi_ds.Raw},
		{true, psi_ds.Gcs},
		{true, psi_ds.BloomFilter},
		{false, psi_ds.Raw},
		{false, psi_ds.Gcs},
		{false, psi_ds.BloomFilter},
	}
	for _, tc := range testCases {
		client, err := psi_client.CreateFromKey(clientKey, tc.revealIntersection)
		if err != nil {
			t.Errorf("Failed to create a PSI client %v", err)
		}
		server, err := psi_server.CreateFromKey(serverKey, tc.revealIntersection)
		if err != nil {
			t.Errorf("Failed to create a PSI server %v", err)
		}

		clientInputs := []string{}
		for i := 0; i < numClientInputs; i++ {
			clientInputs = append(clientInputs, "Element "+string(i))
		}
		serverInputs := []string{}
		for i := 0; i < numServerInputs; i++ {
			serverInputs = append(serverInputs, "Element "+string(i*2))
		}

		// Create the setup
		serverSetup, err := server.CreateSetupMessage(fpr, int64(len(clientInputs)), serverInputs, 0)
		if err != nil {
			t.Errorf("Failed to create serverSetup: %v", err)
		}
		serializedServerSetup, err := proto.Marshal(serverSetup)
		if err != nil {
			t.Errorf("Failed to serialize serverSetup: %v", err)
		}

		serverSetup = &psi_proto.ServerSetup{}
		err = proto.Unmarshal(serializedServerSetup, serverSetup)
		if err != nil {
			t.Errorf("Failed to deserialize serverSetup: %v", err)
		}

		// Create client request
		request, err := client.CreateRequest(clientInputs)
		if err != nil {
			t.Errorf("Failed to create request %v", err)
		}
		serializedRequest, err := proto.Marshal(request)
		if err != nil {
			t.Errorf("Failed to serialize request: %v", err)
		}
		request = &psi_proto.Request{}
		err = proto.Unmarshal(serializedRequest, request)
		if err != nil {
			t.Errorf("Failed to deserialize request: %v", err)
		}

		// Get the response
		response, err := server.ProcessRequest(request)
		if err != nil {
			t.Errorf("Failed to process request: %v", err)
		}
		serializedResponse, err := proto.Marshal(response)
		if err != nil {
			t.Errorf("Failed to serialize response: %v", err)
		}
		response = &psi_proto.Response{}
		err = proto.Unmarshal(serializedResponse, response)
		if err != nil {
			t.Errorf("Failed to deserialize response: %v", err)
		}

		// Compute intersection
		if tc.revealIntersection {
			intersection, err := client.GetIntersection(serverSetup, response)
			if err != nil {
				t.Errorf("failed to compute intersection %v", err)
			}

			if len(intersection) != len(clientInputs)/2 {
				t.Errorf("Intersection size differs! %d", len(intersection))
			}
		} else {
			cardinality, err := client.GetIntersectionSize(serverSetup, response)
			if err != nil {
				t.Errorf("failed to compute intersection %v", err)
			}

			if cardinality != int64(len(clientInputs)/2) {
				t.Errorf("Intersection cardinality differs! %d", cardinality)
			}
		}

		// cleanup
		server.Destroy()
		client.Destroy()
	}
}
