package powerbi

import (
	//"crypto/md5"
	"encoding/base64"
	"github.com/alex-davies/terraform-provider-powerbi/powerbi/internal/api"
	"github.com/hashicorp/terraform/helper/schema"
	"io"
	"strings"
	"time"
)

// ResourcePBIX represents a Power BI PBIX file
func ResourcePBIX() *schema.Resource {
	return &schema.Resource{
		Create: createPBIX,
		Read:   readPBIX,
		Update: updatePBIX,
		Delete: deletePBIX,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"workspace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Workspace in which the PBIX will be added",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the pbix. This will be used as the name for the report and dataset",
			},
			"content_base64": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pbix file as a base64 string",
				// StateFunc: func(val interface{}) string {
				// 	contentBytes, _ := base64.StdEncoding.DecodeString(val.(string))
				// 	hashBytes := md5.Sum(contentBytes)
				// 	hashBase64 := base64.StdEncoding.EncodeToString(hashBytes[:])
				// 	return "md5:" + hashBase64

				// },
			},
			"timeout_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Seconds to wait while publishing the pbix",
				Default:     30,
			},
		},
	}
}

func getContentReader(d *schema.ResourceData) io.Reader {
	contentBase64 := d.Get("content_base64").(string)
	contentReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(contentBase64))

	return contentReader
}

func createPBIX(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	resp, err := client.PostImportInGroup(api.PostImportInGroupRequest{
		GroupID:            d.Get("workspace").(string),
		DatasetDisplayName: d.Get("name").(string),
		NameConflict:       "CreateOrOverwrite",
		Data:               getContentReader(d),
	})
	if err != nil {
		return err
	}

	if client.WaitForImportToSucceed(api.WaitForImportToSucceedRequest{
		ImportID: resp.ID,
		Timeout:  time.Duration(d.Get("timeout_seconds").(int)) * time.Second,
	}) != nil {
		return err
	}

	d.SetId(resp.ID)
	d.Set("workspace", d.Get("workspace").(string))
	d.Set("name", d.Get("name").(string))
	d.Set("content_base64", d.Get("content_base64").(string))
	return nil

}

func readPBIX(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*api.Client)

	return nil
}

func updatePBIX(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*api.Client)

	return nil
}

func deletePBIX(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*api.Client)

	return nil
}
