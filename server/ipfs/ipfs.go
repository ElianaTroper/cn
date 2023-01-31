package ipfs

import (
	"os/exec"
	"github.com/ElianaTroper/cn/server/config"
	"github.com/ElianaTroper/cn/server/daemon"
)

// FUTURE: Check if ipfs nodes are running, if not spin them up
// FUTURE: Run an in house ipfs node if they aren't running and no interface provided

func restart(conf config.CnConf) error {
	// Setting up golang
	// XXX: Assumes docker command is "docker"
	// XXX: Assumes ipfs commands are the default
	// XXX: Assumes ipfs go is running
	command := "ipfs"
	argsShutdown := []string{"shutdown"}
	argsStart := []string{"daemon"}
	if conf.Ipfs.Docker && conf.Ipfs.Docker.IsDocker {
		containerName := conf.Ipfs.Docker.ContainerName
		command = "docker"
		argsShutdown = []string{"exec", containerName, "ipfs shutdown"}
		argsStart = []string{"exec", containerName, "ipfs daemon start"}
	}
	cmdShutdown := exec.Command(command, argsShutdown...)
	// Starting shutdown, will catch it before any other errors though
	err := cmdShutdown.Start()
	if err != nil {
		return err
	}
	cmdStart := exec.Command(command, argsStart...)

	// Setting up js
	jsCommand := "jsipfs"
	jsArgsShutdown := []string{"shutdown"}
	jsArgsStart := []string{"daemon"}
	if conf.Ipfs.Js && conf.Ipfs.Js.Docker && conf.Ipfs.Js.IsDocker {
		containerName := conf.Ipfs.Js.Docker.ContainerName
		jsCommand = "docker"
		jsArgsShutdown = []string{"exec", containerName, "jsipfs shutdown"}
		jsArgsStart := []string{"exec", containerName, "jsipds start"}
	}

	err = cmdShutdown.Wait()
	if err != nil {
		return err
	} // FUTURE: Add a timeout
	err = cmdStart.Start()
	if err != nil {
		return err
	} // FUTURE: Watch this process for later errors and pass up
	if conf.Ipfs.Js {
		jsCmdShutdown := exec.Command(jsCommand, jsArgsShutdown...)
		err = jsCmdShutdown.Run()
		if err != nil {
			return err
		} // FUTURE: Add a timeout
		jsCmdStart := exec.Command(jsCommand, jsArgsStart...)
		err = jsCmdStart.Start()
		if err != nil {
			return err
		} // FUTURE: Watch this process for later errors and pass up
	}
	return nil
}

func Init(conf config.CnConf) error {
	// FUTURE: Check that this has not already been done
	ipfsSettings := conf.Ipfs
	writeIpfsConf := false
	var ipfsConf map[string]interface{}
	ipfsConfPath := ipfsSettings.Conf.Path
	err := config.LoadUnstructured(ipfsConfPath, &ipfsConf)
	if err != nil {
		return err
	}
	// If overrides are set:
	overridesPath := ipfsSettings.Conf.Overrides
	if overridesPath != "" {
		writeIpfsConf = true
		// Load the overrides
		var loadedOverrides map[string]interface{}
		err = config.LoadUnstructured(overridesPath, &loadedOverrides)
		if err != nil {
			return err
		}
		// Replace the values
		ipfsConf, err = config.ReplaceOrAdd(loadedOverrides, ipfsConf)
		if err != nil {
			return err
		}
	}
	// If we're linking a go and js ipfs node:
	jsConf := ipfsSettings.Js
	if jsConf != nil && jsConf.EnableLink {
		writeIpfsConf = true
		// TODO
	}
	// If we're inferring announce addresses:
	if ipfsSettings.InferAnnounce {
		writeIpfsConf = true
		// TODO
	}

	// If changes have occurred, replace the conf
	if writeIpfsConf {
		err = config.WriteUnstructured(ipfsConf)
		if err != nil {
			return err
		}
		// FUTURE: Write out the old config
		return restart(conf)
	}
	return nil
}

// TODO: Also need to do js settings
