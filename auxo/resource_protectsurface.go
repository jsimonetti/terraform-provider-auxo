package auxo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/on2itsecurity/go-auxo/zerotrust"
)

func resourceProtectSurface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProtectSurfaceCreate,
		ReadContext:   resourceProtectSurfaceRead,
		UpdateContext: resourceProtectSurfaceUpdate,
		DeleteContext: resourceProtectSurfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the resource/segment",
			},
			"uniqueness_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "Unique key to generate the ID - only needed for parallel import",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the segment",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the segment",
			},
			"main_contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of main contact in text)",
			},
			"security_contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of security contact in text",
			},
			"in_control_boundary": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Is this protect surface in the control boundary (your responsibility)",
			},
			"in_zero_trust_focus": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Is this protect surface in the zero trust focus (actively maintained and monitored)",
			},
			"relevance": {
				Type:        schema.TypeInt,
				Default:     60,
				Optional:    true,
				Description: "Relevance 0-100 of the segment",
			},
			"confidentiality": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Confidentiality score (number 1-5)",
			},
			"integrity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Integrity score (number 1-5)",
			},
			"availability": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Availability score (number 1-5)",
			},
			"data_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Contains data tags, defining the data residing in the protect surface",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"compliance_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Contains compliance tags, defining compliancy requirements of the protect surface",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"customer_labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Contains customer labels in Key-Value-Pair format",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"soc_tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Contains tags, which are used by the SOC engineers",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"maturity_step1": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "maturity step 1 - defining the protect surface",
			},
			"maturity_step2": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "maturity step 2 - map the transaction flows",
			},
			"maturity_step3": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "maturity step 3 - architect your environment",
			},
			"maturity_step4": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "maturity step 4 - zero trust policy",
			},
			"maturity_step5": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "maturity step 5 - monitor and maintain",
			},
		},
	}
}

func resourceProtectSurfaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var ps = new(zerotrust.ProtectSurface)
	ps.UniquenessKey = d.Get("uniqueness_key").(string)
	ps.Name = d.Get("name").(string)
	ps.Description = d.Get("description").(string)
	ps.MainContactPersonID = d.Get("main_contact").(string)
	ps.SecurityContactPersonID = d.Get("security_contact").(string)
	ps.InControlBoundary = d.Get("in_control_boundary").(bool)
	ps.InZeroTrustFocus = d.Get("in_zero_trust_focus").(bool)
	ps.Relevance = d.Get("relevance").(int)
	ps.Confidentiality = d.Get("confidentiality").(int)
	ps.Integrity = d.Get("integrity").(int)
	ps.Availability = d.Get("availability").(int)
	ps.DataTags = createStringSliceFromListInput(d.Get("data_tags").(*schema.Set).List())
	ps.ComplianceTags = createStringSliceFromListInput(d.Get("compliance_tags").(*schema.Set).List())
	ps.SocTags = createStringSliceFromListInput(d.Get("soc_tags").(*schema.Set).List())
	cl := make(map[string]string)
	for k, v := range d.Get("customer_labels").(map[string]any) {
		cl[k] = v.(string)
	}
	ps.CustomerLabels = cl
	ps.Maturity.Step1 = d.Get("maturity_step1").(int)
	ps.Maturity.Step2 = d.Get("maturity_step2").(int)
	ps.Maturity.Step3 = d.Get("maturity_step3").(int)
	ps.Maturity.Step4 = d.Get("maturity_step4").(int)
	ps.Maturity.Step5 = d.Get("maturity_step5").(int)

	result, err := apiClient.ZeroTrust.CreateProtectSurfaceByObject(*ps, false)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)

	resourceProtectSurfaceRead(ctx, d, m)

	return diags
}

func resourceProtectSurfaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	ps, err := apiClient.ZeroTrust.GetProtectSurfaceByID(d.Id())

	if err != nil {
		apiError := getAPIError(err)

		//NotExists
		if apiError.ID == "410" {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("id", ps.ID)
	d.Set("uniqueness_key", ps.UniquenessKey)
	d.Set("name", ps.Name)
	d.Set("description", ps.Description)
	d.Set("main_contact", ps.MainContactPersonID)
	d.Set("security_contact", ps.SecurityContactPersonID)
	d.Set("in_control_boundary", ps.InControlBoundary)
	d.Set("in_zero_trust_focus", ps.InZeroTrustFocus)
	d.Set("relevance", ps.Relevance)
	d.Set("confidentiality", ps.Confidentiality)
	d.Set("integrity", ps.Integrity)
	d.Set("availability", ps.Availability)
	d.Set("data_tags", ps.DataTags)
	d.Set("compliance_tags", ps.ComplianceTags)
	d.Set("customer_labels", ps.CustomerLabels)
	d.Set("soc_tags", ps.SocTags)
	d.Set("maturity_step1", ps.Maturity.Step1)
	d.Set("maturity_step2", ps.Maturity.Step2)
	d.Set("maturity_step3", ps.Maturity.Step3)
	d.Set("maturity_step4", ps.Maturity.Step4)
	d.Set("maturity_step5", ps.Maturity.Step5)

	return diags
}

func resourceProtectSurfaceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	ps, err := apiClient.ZeroTrust.GetProtectSurfaceByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("uniqueness_key") {
		ps.UniquenessKey = d.Get("uniqueness_key").(string)
	}
	if d.HasChange("name") {
		ps.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		ps.Description = d.Get("description").(string)
	}
	if d.HasChange("main_contact") {
		ps.MainContactPersonID = d.Get("main_contact").(string)
	}
	if d.HasChange("security_contact") {
		ps.SecurityContactPersonID = d.Get("security_contact").(string)
	}
	if d.HasChange("in_control_boundary") {
		ps.InControlBoundary = d.Get("in_control_boundary").(bool)
	}
	if d.HasChange("in_zero_trust_focus") {
		ps.InZeroTrustFocus = d.Get("in_zero_trust_focus").(bool)
	}
	if d.HasChange("relevance") {
		ps.Relevance = d.Get("relevance").(int)
	}
	if d.HasChange("confidentiality") {
		ps.Confidentiality = d.Get("confidentiality").(int)
	}
	if d.HasChange("integrity") {
		ps.Integrity = d.Get("integrity").(int)
	}
	if d.HasChange("availability") {
		ps.Availability = d.Get("availability").(int)
	}
	if d.HasChange("data_tags") {
		ps.DataTags = createStringSliceFromListInput(d.Get("data_tags").(*schema.Set).List())
	}
	if d.HasChange("compliance_tags") {
		ps.ComplianceTags = createStringSliceFromListInput(d.Get("compliance_tags").(*schema.Set).List())
	}
	if d.HasChange("customer_labels") {
		cl := make(map[string]string)
		for k, v := range d.Get("customer_labels").(map[string]any) {
			cl[k] = v.(string)
		}
		ps.CustomerLabels = cl
	}
	if d.HasChange("soc_tags") {
		ps.SocTags = createStringSliceFromListInput(d.Get("soc_tags").(*schema.Set).List())
	}
	if d.HasChange("maturity_step1") {
		ps.Maturity.Step1 = d.Get("maturity_step1").(int)
	}
	if d.HasChange("maturity_step2") {
		ps.Maturity.Step2 = d.Get("maturity_step2").(int)
	}
	if d.HasChange("maturity_step3") {
		ps.Maturity.Step3 = d.Get("maturity_step3").(int)
	}
	if d.HasChange("maturity_step4") {
		ps.Maturity.Step4 = d.Get("maturity_step4").(int)
	}
	if d.HasChange("maturity_step5") {
		ps.Maturity.Step5 = d.Get("maturity_step5").(int)
	}

	_, err = apiClient.ZeroTrust.UpdateProtectSurface(*ps)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProtectSurfaceRead(ctx, d, m)
}

func resourceProtectSurfaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	provider := m.(*AuxoProvider)
	apiClient := provider.APIClient

	var diags diag.Diagnostics

	err := apiClient.ZeroTrust.DeleteProtectSurfaceByID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func createStringSliceFromListInput(inputList []interface{}) []string {
	output := make([]string, len(inputList))
	for k, v := range inputList {
		output[k] = v.(string)
	}

	return output
}
