package main

import "k8s.io/client-go/tools/leaderelection"

func main()  {
	election, _ := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{

	})
	election.IsLeader()
}
