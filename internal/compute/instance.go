package compute

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"google.golang.org/api/compute/v1"
	"strings"
)

type Instance struct {
	Instance *compute.Instance
}

func (i *Instance) GetOutput() (out string) {
	instance := i.Instance

	instanceName := instance.Name
	powerState := instance.Status
	sizeInfo := strings.Split(instance.MachineType, "/")
	instanceSize := sizeInfo[len(sizeInfo)-1]

	out = fmt.Sprintf(`"%s" powerstate=%s size=%s`, instanceName, powerState, instanceSize)

	return
}

func (i *Instance) GetStatus() int {
	var state int

	// The status of the instance, possible values include: "DEPROVISIONING", "PROVISIONING", "REPAIRING", "STAGING",
	//"STOPPED", "STOPPING", "SUSPENDED", "SUSPENDING", "TERMINATED", "RUNNING"
	switch i.Instance.Status {
	case "DEPROVISIONING", "PROVISIONING", "REPAIRING", "STAGING", "STOPPED", "STOPPING", "SUSPENDED", "SUSPENDING", "TERMINATED":
		state = check.Critical
	case "RUNNING":
		state = check.OK
	default:
		state = check.Critical
	}

	return state
}
