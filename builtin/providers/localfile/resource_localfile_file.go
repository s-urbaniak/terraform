package localfile

import (
	"fmt"
	"os"

	"io/ioutil"

	"crypto/sha1"
	"encoding/hex"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalFileCreate,
		Read:   resourceLocalFileRead,
		Delete: resourceLocalFileDelete,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: diffSuppress,
			},
			"destination": {
				Type:        schema.TypeString,
				Description: "Path to the output file",
				Required:    true,
				ForceNew:    true,
			},
			"sensitive": {
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
			},
		},
	}
}

func diffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if data, ok := d.GetOk("sensitive"); ok {
		fmt.Printf("!!!!! value %v type %T\n", data, data)
		return true
	}

	return false
}

func resourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	// If the output file doesn't exist, mark the resource for creation.
	outputPath := d.Get("destination").(string)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceLocalFileCreate(d *schema.ResourceData, _ interface{}) error {
	content := d.Get("content").(string)
	destination := d.Get("destination").(string)

	destinationDir := path.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		if err := os.MkdirAll(destinationDir, 0777); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(destination, []byte(content), 0777); err != nil {
		return err
	}

	checksum := sha1.Sum([]byte(content))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}

func resourceLocalFileDelete(d *schema.ResourceData, _ interface{}) error {
	os.Remove(d.Get("destination").(string))
	return nil
}
