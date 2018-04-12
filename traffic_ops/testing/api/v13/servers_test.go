package v13

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"testing"
	"github.com/apache/incubator-trafficcontrol/lib/go-log"
	"github.com/apache/incubator-trafficcontrol/lib/go-tc"
)

func TestServers(t *testing.T) {

	CreateTestCDNs(t)
	CreateTestTypes(t)
	CreateTestProfiles(t)
	CreateTestStatuses(t)
	CreateTestDivisions(t)
	CreateTestRegions(t)
	CreateTestPhysLocations(t)
	CreateTestCacheGroups(t)
	CreateTestServers(t)

	UpdateTestServers(t)
	GetTestServers(t)

	DeleteTestServers(t)
	DeleteTestCacheGroups(t)
	DeleteTestPhysLocations(t)
	DeleteTestRegions(t)
	DeleteTestDivisions(t)
	DeleteTestStatuses(t)
	DeleteTestProfiles(t)
	DeleteTestTypes(t)
	DeleteTestCDNs(t)

}

func CreateTestServers(t *testing.T) {

	// GET EDGE1 profile
	resp, _, err := TOSession.GetProfileByName("EDGE1")
	if err != nil {
		t.Errorf("cannot GET Profiles - %v\n", err)
	}
	respProfile := resp[0]

	// GET EDGE type
	resp2, _, err := TOSession.GetTypeByName("EDGE")
	if err != nil {
		t.Errorf("cannot GET Division by name: EDGE - %v\n", err)
	}
	respType := resp2[0]

	// GET ONLINE status
	resp3, _, err := TOSession.GetStatusByName("ONLINE")
	if err != nil {
		t.Errorf("cannot GET Status by name: ONLINE - %v\n", err)
	}
	respStatus := resp3[0]

	// GET Denver physlocation
	resp4, _, err := TOSession.GetPhysLocationByName("Denver")
	if err != nil {
		t.Errorf("cannot GET PhysLocation by name: Denver - %v\n", err)
	}
	respPhysLocation := resp4[0]

	// GET cachegroup1 cachegroup
	resp5, _, err := TOSession.GetCacheGroupByName("cachegroup1")
	if err != nil {
		t.Errorf("cannot GET CacheGroup by name: cachegroup1 - %v\n", err)
	}
	respCacheGroup := resp5[0]

	// loop through servers, assign FKs and create
	for _, server := range testData.Servers {
		server.CDNID = respProfile.CDNID
		server.ProfileID = respProfile.ID
		server.TypeID = respType.ID
		server.StatusID = respStatus.ID
		server.PhysLocationID = respPhysLocation.ID
		server.CachegroupID = respCacheGroup.ID

		resp, _, err := TOSession.CreateServer(server)
		log.Debugln("Response: ", server.HostName, " ", resp)
		if err != nil {
			t.Errorf("could not CREATE servers: %v\n", err)
		}
	}

}

func GetTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by name: %v - %v\n", err, resp)
		}
	}
}

func UpdateTestServers(t *testing.T) {

	firstServer := testData.Servers[0]
	// Retrieve the Server by hostname so we can get the id for the Update
	resp, _, err := TOSession.GetServerByHostName(firstServer.HostName)
	if err != nil {
		t.Errorf("cannot GET Server by hostname: %v - %v\n", firstServer.HostName, err)
	}
	remoteServer := resp[0]
	updatedServerInterface := "bond1"
	updatedServerRack := "RR 119.03"

	// update rack and interfaceName values on server
	remoteServer.InterfaceName = updatedServerInterface
	remoteServer.Rack = updatedServerRack
	var alert tc.Alerts
	alert, _, err = TOSession.UpdateServerByID(remoteServer.ID, remoteServer)
	if err != nil {
		t.Errorf("cannot UPDATE Server by hostname: %v - %v\n", err, alert)
	}

	// Retrieve the Profile to check Profile name got updated
	resp, _, err = TOSession.GetServerByID(remoteServer.ID)
	if err != nil {
		t.Errorf("cannot GET Server by ID: %v - %v\n", remoteServer.HostName, err)
	}

	respServer := resp[0]
	if respServer.InterfaceName != updatedServerInterface || respServer.Rack != updatedServerRack {
		t.Errorf("results do not match actual: %s, expected: %s\n", respServer.InterfaceName, updatedServerInterface)
		t.Errorf("results do not match actual: %s, expected: %s\n", respServer.Rack, updatedServerRack)
	}

}

func DeleteTestServers(t *testing.T) {

	for _, server := range testData.Servers {
		resp, _, err := TOSession.GetServerByHostName(server.HostName)
		if err != nil {
			t.Errorf("cannot GET Server by hostname: %v - %v\n", server.HostName, err)
		}
		if len(resp) > 0 {
			respServer := resp[0]

			delResp, _, err := TOSession.DeleteServerByID(respServer.ID)
			if err != nil {
				t.Errorf("cannot DELETE Server by hostname: %v - %v\n", err, delResp)
			}

			// Retrieve the Server to see if it got deleted
			serv, _, err := TOSession.GetServerByHostName(server.HostName)
			if err != nil {
				t.Errorf("error deleting Server hostname: %s\n", err.Error())
			}
			if len(serv) > 0 {
				t.Errorf("expected Server hostname: %s to be deleted\n", server.HostName)
			}
		}
	}
}