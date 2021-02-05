package azureBilling

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/adeturner/observability"
)

func (vml *vmLookup) print(cnt int) {

	i := 0

	for k, v := range vml.items {
		observability.Logger("Info", fmt.Sprintf("%s -> %v\n", k, v))
		i++
		if i > cnt {
			break
		}
	}
}

func (vml *vmLookup) printCount() {
	observability.Logger("Info", fmt.Sprintf("vmSizeLookup has %d records\n", len(vml.items)))
}

func (vml *vmLookup) init() {
	vml.items = make(map[string]vmLookupItem)
}

func (vml *vmLookup) Read(fileLocation string) error {

	vml.init()

	f, err := os.Open(fileLocation)
	if err != nil {
		observability.Logger("Error", fmt.Sprintf("Unable to read input file=%s err=%s", fileLocation, err))
	}
	defer f.Close()

	cnt := 0

	if err == nil {

		r := csv.NewReader(f)
		for {

			record, err := r.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				observability.Logger("Error", fmt.Sprintf("Unable to parse file as CSV; file=%s err=%s", fileLocation, err))
				break
			}

			cnt++

			// skip the first row (header)
			if cnt > 1 {
				i := vmLookupItem{}
				i.setValues(record)

				vml.items[strings.ToLower(i.VM)] = i
			}

		}
	}

	observability.LogMemory("Info")
	vml.printCount()
	// vml.print(10)

	return err

}

func (vml *vmLookup) get(vmName string) (vmLookupItem, bool) {

	key := fmt.Sprintf("%s", strings.ToLower(vmName))

	vmli, ok := vml.items[key]
	if !ok {
		observability.Logger("Error", fmt.Sprintf("Unable to find reportingCategoryLookupItem for key=%s", key))
	}

	return vmli, ok
}