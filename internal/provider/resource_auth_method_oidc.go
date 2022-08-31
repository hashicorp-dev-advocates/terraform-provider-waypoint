package waypoint

import (
	"context"
	"fmt"
	"github.com/hashicorp-dev-advocates/waypoint-client/pkg/client"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAuthMethodOidc() *schema.Resource {
	return &schema.Resource{
		Description: "Auth method OIDC resource manages OIDC auth methods in Waypoint.",

		CreateContext: resourceAuthMethodOidcCreate,
		ReadContext:   resourceAuthMethodOidcRead,
		UpdateContext: resourceAuthMethodOidcCreate,
		DeleteContext: resourceAuthMethodOidcDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of OIDC auth method",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Friendly display name of OIDC auth method",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of auth method",
			},
			"accessor_selector": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Client ID of OIDC provider",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "client secret for OIDC provider",
				Sensitive:   true,
			},
			"discovery_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Discovery URL for OIDC provider",
			},
			"allowed_redirect_urls": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Allowed URI for auth redirection.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"claim_mappings": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Mapping of a claim to a variable value for the access selector",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"list_claim_mappings": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Same as claim-mapping but for list values",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"discovery_ca_pem": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Optional CA certificate chain to validate the discovery URL. Multiple CA certificates can be specified to support easier rotation",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"signing_algs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The signing algorithms supported by the OIDC connect server. If this isn't specified, this will default to RS256 since that should be supported according to the RFC. The string values here should be valid OIDC signing algorithms",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"scopes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The optional claims scope requested.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auds": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The optional audience claims required",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceAuthMethodOidcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	wp := m.(*WaypointClient).conn

	authMethodConfig := client.DefaultAuthMethodConfig()

	if name, ok := d.Get("name").(string); ok {
		authMethodConfig.Name = name
	}

	if displayName, ok := d.Get("display_name").(string); ok {
		authMethodConfig.DisplayName = displayName
	}

	if description, ok := d.Get("description").(string); ok {
		authMethodConfig.Description = description
	}

	if accessorSelector, ok := d.Get("accessor_selector").(string); ok {
		authMethodConfig.AccessSelector = accessorSelector
	}

	oidcConfig := client.DefaultOidcConfig()
	oidcConfig.ClientId = d.Get("client_id").(string)
	oidcConfig.DiscoveryUrl = d.Get("discovery_url").(string)

	if clientSecret, ok := d.Get("client_secret").(string); ok {
		oidcConfig.ClientSecret = clientSecret
	}

	redirectUrls := d.Get("allowed_redirect_urls").([]interface{})
	strRedirectUrl := make([]string, len(redirectUrls))

	for i, v := range redirectUrls {
		strRedirectUrl[i] = fmt.Sprint(v)
	}

	oidcConfig.AllowedRedirectUris = strRedirectUrl
	if claimMappings, ok := d.Get("claim_mappings").(map[string]interface{}); ok {
		cMappings := make(map[string]string)

		for k, v := range claimMappings {
			strKey := fmt.Sprintf("%v", k)
			strValue := fmt.Sprintf("%v", v)
			cMappings[strKey] = strValue
		}

		oidcConfig.ClaimMappings = cMappings
	}

	if listClaimMappings, ok := d.Get("list_claim_mappings").(map[string]interface{}); ok {
		mappings := make(map[string]string)

		for k, v := range listClaimMappings {
			strKey := fmt.Sprintf("%v", k)
			strValue := fmt.Sprintf("%v", v)
			mappings[strKey] = strValue
		}

		oidcConfig.ListClaimMappings = mappings
	}

	if audience, ok := d.Get("auds").([]interface{}); ok {
		audSlice := make([]string, len(audience))

		for i, v := range audience {
			audSlice[i] = fmt.Sprint(v)
		}

		oidcConfig.Auds = audSlice
	}

	if scopes, ok := d.Get("scopes").([]interface{}); ok {
		scopeSlice := make([]string, len(scopes))

		for i, v := range scopes {
			scopeSlice[i] = fmt.Sprint(v)
		}

		oidcConfig.Scopes = scopeSlice
	}

	if signingAlgos, ok := d.Get("signing_algs").([]interface{}); ok {
		//signingAlgs := d.Get("signing_algs").([]interface{})
		signingAlgsList := make([]string, len(signingAlgos))

		for i, value := range signingAlgos {
			signingAlgsList[i] = fmt.Sprint(value)
		}

		oidcConfig.SigningAlgs = signingAlgsList

	}

	if discoveryCaPem, ok := d.Get("discovery_ca_pem").([]interface{}); ok {
		discoveryCaPemSlice := make([]string, len(discoveryCaPem))

		for i, v := range discoveryCaPem {
			discoveryCaPemSlice[i] = fmt.Sprint(v)
		}

		oidcConfig.DiscoveryCaPem = discoveryCaPemSlice
	}

	authMethod, err := wp.UpsertOidc(context.TODO(), oidcConfig, authMethodConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(authMethod.Name)

	tflog.Trace(ctx, "created a resource")
	tflog.Debug(ctx, fmt.Sprintf("%s", oidcConfig))

	return resourceAuthMethodOidcRead(ctx, d, m)
}

func resourceAuthMethodOidcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	wp := m.(*WaypointClient).conn

	name := d.Get("name").(string)
	am, err := wp.GetOidcAuthMethod(context.TODO(), name)
	if err != nil {
		return diag.FromErr(err)
	}

	//am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.SigningAlgs
	d.SetId(name)
	d.Set("display_name", am.AuthMethod.DisplayName)
	d.Set("description", am.AuthMethod.Description)
	d.Set("accessor_selector", am.AuthMethod.AccessSelector)
	d.Set("client_id", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.ClientId)
	d.Set("client_secret", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.ClientSecret)
	d.Set("discovery_url", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.DiscoveryUrl)
	d.Set("allowed_redirect_urls", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.AllowedRedirectUris)
	d.Set("claim_mappings", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.ClaimMappings)
	d.Set("list_claim_mappings", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.ListClaimMappings)
	d.Set("discovery_ca_pem", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.DiscoveryCaPem)
	d.Set("signing_algs", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.SigningAlgs)
	d.Set("scopes", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.Scopes)
	d.Set("auds", am.AuthMethod.Method.(*gen.AuthMethod_Oidc).Oidc.Auds)

	return nil
}

func resourceAuthMethodOidcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	wp := m.(*WaypointClient).conn

	err := wp.DeleteOidc(context.TODO(), d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
