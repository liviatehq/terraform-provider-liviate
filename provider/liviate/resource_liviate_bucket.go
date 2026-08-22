package liviate

import (
	"fmt"
	"log"
	"strings"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudStackBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudStackBucketCreate,
		Read:   resourceCloudStackBucketRead,
		Update: resourceCloudStackBucketUpdate,
		Delete: resourceCloudStackBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"objectstorageid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"quota": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encryption": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"versioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"objectlocking": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"account": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"objectstore_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"accesskey": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"usersecretkey": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCloudStackBucketCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating bucket %s", d.Get("name").(string))

	cs := meta.(*cloudstack.CloudStackClient)

	name := d.Get("name").(string)
	objectstorageid := d.Get("objectstorageid").(string)
	quota := d.Get("quota").(int)

	// If no object storage pool is specified, try to discover it automatically
	if objectstorageid == "" {
		l, err := cs.StoragePool.ListObjectStoragePools(cs.StoragePool.NewListObjectStoragePoolsParams())
		if err != nil {
			return fmt.Errorf("error listing object storage pools: %s", err)
		}
		switch len(l.ObjectStoragePools) {
		case 0:
			return fmt.Errorf("no object storage pool found; configure one or set objectstorageid")
		case 1:
			objectstorageid = l.ObjectStoragePools[0].Id
		default:
			return fmt.Errorf("multiple object storage pools found; set objectstorageid explicitly")
		}
	}

	p := cs.ObjectStore.NewCreateBucketParams(name, objectstorageid, quota)

	if v, ok := d.GetOk("policy"); ok {
		p.SetPolicy(v.(string))
	}
	if v, ok := d.GetOk("encryption"); ok {
		p.SetEncryption(v.(bool))
	}
	if v, ok := d.GetOk("versioning"); ok {
		p.SetVersioning(v.(bool))
	}
	if v, ok := d.GetOk("objectlocking"); ok {
		p.SetObjectlocking(v.(bool))
	}
	if v, ok := d.GetOk("account"); ok {
		p.SetAccount(v.(string))
	}
	if v, ok := d.GetOk("domain_id"); ok {
		p.SetDomainid(v.(string))
	}
	if v, ok := d.GetOk("project_id"); ok {
		p.SetProjectid(v.(string))
	}

	r, err := cs.ObjectStore.CreateBucket(p)
	if err != nil {
		return fmt.Errorf("error creating bucket %s: %s", name, err)
	}

	d.SetId(r.Id)

	return resourceCloudStackBucketRead(d, meta)
}

func resourceCloudStackBucketRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading bucket %s", d.Id())

	cs := meta.(*cloudstack.CloudStackClient)

	id := d.Id()

	bucket, count, err := cs.ObjectStore.GetBucketByID(id)
	if err != nil {
		if count == 0 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error getting bucket %s: %s", id, err)
	}

	d.Set("name", bucket.Name)
	d.Set("objectstorageid", bucket.Objectstorageid)
	d.Set("quota", bucket.Quota)
	d.Set("policy", bucket.Policy)
	d.Set("encryption", bucket.Encryption)
	d.Set("versioning", bucket.Versioning)
	d.Set("objectlocking", bucket.Objectlocking)
	d.Set("url", bucket.Url)
	d.Set("size", bucket.Size)
	d.Set("state", bucket.State)
	d.Set("objectstore_provider", bucket.Provider)
	d.Set("accesskey", bucket.Accesskey)
	d.Set("usersecretkey", bucket.Usersecretkey)
	d.Set("account", bucket.Account)
	d.Set("domain_id", bucket.Domainid)
	d.Set("project_id", bucket.Projectid)

	return nil
}

func resourceCloudStackBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating bucket %s", d.Id())

	cs := meta.(*cloudstack.CloudStackClient)

	id := d.Id()
	p := cs.ObjectStore.NewUpdateBucketParams(id)

	if d.HasChange("quota") {
		p.SetQuota(d.Get("quota").(int))
	}

	if d.HasChange("policy") {
		p.SetPolicy(d.Get("policy").(string))
	}

	if d.HasChange("encryption") {
		p.SetEncryption(d.Get("encryption").(bool))
	}

	if d.HasChange("versioning") {
		p.SetVersioning(d.Get("versioning").(bool))
	}

	if d.HasChanges("name", "objectstorageid") {
		return fmt.Errorf("cannot update immutable field name or objectstorageid on bucket %s", id)
	}

	_, err := cs.ObjectStore.UpdateBucket(p)
	if err != nil {
		return fmt.Errorf("error updating bucket %s: %s", id, err)
	}

	return resourceCloudStackBucketRead(d, meta)
}

func resourceCloudStackBucketDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting bucket %s", d.Id())

	cs := meta.(*cloudstack.CloudStackClient)

	id := d.Id()
	p := cs.ObjectStore.NewDeleteBucketParams(id)

	_, err := cs.ObjectStore.DeleteBucket(p)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf(
			"InvalidParameterValueException: Can't find bucket by id '%s'", id)) {
			return nil
		}
		return fmt.Errorf("error deleting bucket %s: %s", id, err)
	}

	return nil
}
