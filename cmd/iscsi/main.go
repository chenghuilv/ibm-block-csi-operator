package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/IBM/ibm-block-csi-driver-operator/pkg/iscsi"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("cmd")

const (
	LOGIN  = "login"
	LOGOUT = "logout"
	WATCH  = "watch"
)

func main() {
	logf.SetLogger(zap.Logger())

	loginCmd := flag.NewFlagSet(LOGIN, flag.ExitOnError)
	loginPortals := loginCmd.String("portals", "", "portals")

	logoutCmd := flag.NewFlagSet(LOGOUT, flag.ExitOnError)
	logoutPortals := logoutCmd.String("portals", "", "portals")

	watchCmd := flag.NewFlagSet(WATCH, flag.ExitOnError)

	if len(os.Args) < 2 {
		log.Info("expected a subcommand", "candidate", []string{LOGIN, LOGOUT, WATCH})
		os.Exit(1)
	}

	switch os.Args[1] {

	case LOGIN:
		loginCmd.Parse(os.Args[2:])
		login(*loginPortals)
	case LOGOUT:
		logoutCmd.Parse(os.Args[2:])
		logout(*logoutPortals)
	case WATCH:
		watchCmd.Parse(os.Args[2:])
		watch()
	default:
		log.Info("expected a subcommand", "candidate", []string{LOGIN, LOGOUT, WATCH})
		os.Exit(1)
	}
}

func login(portals string) {

	if portals == "" {
		log.Info("--portals is required!")
		os.Exit(1)
	}
	log.Info("Starting to login portals: " + portals)

	for _, portal := range strings.Split(portals, ",") {
		err := iscsi.DiscoverAndLogin(portal)
		if err != nil {
			log.Error(err, "Failed to login portal "+portal)
		}
	}
	log.Info("Finished to login portals")
}

func logout(portals string) {

	if portals == "" {
		log.Info("--portals is required!")
		os.Exit(1)
	}
	log.Info("Starting to logout portals: " + portals)

	for _, portal := range strings.Split(portals, ",") {
		err := iscsi.DiscoverAndLogout(portal)
		if err != nil {
			log.Error(err, "Failed to logout portal "+portal)
		}
	}
	log.Info("Finished to logout portals")
}

func watch() {
	for {
		time.Sleep(time.Hour * 24)
	}
}