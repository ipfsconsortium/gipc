package commands

import (
	"errors"
	"fmt"
	"os"

	cfg "github.com/ipfsconsortium/gipc/config"
	ipfsclient "github.com/ipfsconsortium/gipc/ipfsc"
	"github.com/ipfsconsortium/gipc/service"
	sto "github.com/ipfsconsortium/gipc/storage"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	errInvalidParameters = errors.New("invalid parameters")
)

// DumpDb command
func DumpDb(cmd *cobra.Command, args []string) {

	must(loadStorage())

	storage.Dump(os.Stdout)
}

// InitDb command
func InitDb(cmd *cobra.Command, args []string) {

	must(loadStorage())

	storage.SetGlobals(sto.GlobalsEntry{
		CurrentQuota: 0,
	})

}

func IpfscInfo(cmd *cobra.Command, args []string) {

	must(loadEthClients())
	must(loadIPFSC())

	info, err := ipfsc.ENS().Info(cfg.C.EnsNames.Local)
	if err != nil {
		log.WithError(err).Error("Failed to get info")
		return
	}
	manifest, err := ipfsc.Read(cfg.C.EnsNames.Local)
	if err != nil {
		log.WithError(err).Error("Failed to read manifest")
		fmt.Println(info)
		return
	}
	pinningManifest := manifest.(*ipfsclient.PinningManifest)
	for _, ipfshash := range pinningManifest.Pin {
		info += "\nPin: " + ipfshash
	}
	fmt.Println(info)
}

func IpfscInit(cmd *cobra.Command, args []string) {

	must(loadEthClients())
	must(loadIPFSC())

	quotum := args[0]

	var manifest ipfsclient.PinningManifest
	manifest.Quotum = quotum

	if err := ipfsc.Write(cfg.C.EnsNames.Local, &manifest); err != nil {
		log.Error("Failed to init ", err)
		return
	}
	log.Info("Sucessfully initialized")

}

func IpfscAdd(cmd *cobra.Command, args []string) {

	must(loadEthClients())
	must(loadIPFSC())

	m, err := ipfsc.Read(cfg.C.EnsNames.Local)
	if err != nil {
		log.Error("Failed to read manifest", err)
		return
	}

	manifest := m.(*ipfsclient.PinningManifest)
	for _, ipfshash := range args {
		manifest.Pin = append(manifest.Pin, ipfshash)
	}

	if err := ipfsc.Write(cfg.C.EnsNames.Local, manifest); err != nil {
		log.Error("Failed to write manifest ", err)
		return
	}
	log.Info("Manifest sucessfully updated")
}

func IpfscRemove(cmd *cobra.Command, args []string) {

	must(loadEthClients())
	must(loadIPFSC())

	m, err := ipfsc.Read(cfg.C.EnsNames.Local)
	if err != nil {
		log.Error("Failed to read manifest", err)
		return
	}

	remove := make(map[string]bool)
	for _, ipfshash := range args {
		remove[ipfshash] = true
	}

	manifest := m.(*ipfsclient.PinningManifest)
	for i, ipfshash := range args {
		if _, ok := remove[ipfshash]; !ok {
			manifest.Pin = append(manifest.Pin[:i], manifest.Pin[i+1:]...)
		}
	}

	if err := ipfsc.Write(cfg.C.EnsNames.Local, manifest); err != nil {
		log.Error("Failed to write manifest ", err)
		return
	}
	log.Info("Manifest sucessfully updated")
}

// Serve command
func Serve(cmd *cobra.Command, args []string) {

	must(load(true))

	service.NewService(
		ipfsc, storage,
	).Serve()

}
