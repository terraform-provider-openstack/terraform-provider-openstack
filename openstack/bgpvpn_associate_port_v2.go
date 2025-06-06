package openstack

import (
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandBGPVPNPortAssociateRoutesV2(routes []any) []bgpvpns.PortRoutes {
	res := make([]bgpvpns.PortRoutes, len(routes))
	for i, r := range routes {
		res[i] = expandBGPVPNPortAssociateRouteV2(r.(map[string]any), false)
	}

	return res
}

func expandBGPVPNPortAssociateRoutesUpdateV2(d *schema.ResourceData) []bgpvpns.PortRoutes {
	oldRoutesRaw, newRoutesRaw := d.GetChange("routes")
	oldRoutes, newRoutes := oldRoutesRaw.(*schema.Set).List(), newRoutesRaw.(*schema.Set).List()
	res := make([]bgpvpns.PortRoutes, len(newRoutes))

	for i, nr := range newRoutes {
		or := oldRoutes[i].(map[string]any)
		nr := nr.(map[string]any)
		olp := or["local_pref"].(int)
		nlp := nr["local_pref"].(int)
		// set local_pref to 0 only, when it was set before and is not set now
		res[i] = expandBGPVPNPortAssociateRouteV2(nr, olp > 0 && nlp == 0)
	}

	return res
}

func expandBGPVPNPortAssociateRouteV2(route map[string]any, enforceLocalPref bool) bgpvpns.PortRoutes {
	res := bgpvpns.PortRoutes{
		Type:     route["type"].(string),
		Prefix:   route["prefix"].(string),
		BGPVPNID: route["bgpvpn_id"].(string),
	}
	if v := route["local_pref"].(int); v > 0 || enforceLocalPref {
		res.LocalPref = &v
	}

	return res
}

func flattenBGPVPNPortAssociateRoutesV2(routes []bgpvpns.PortRoutes) []map[string]any {
	res := make([]map[string]any, len(routes))
	for i, r := range routes {
		res[i] = map[string]any{
			"type":       r.Type,
			"prefix":     r.Prefix,
			"bgpvpn_id":  r.BGPVPNID,
			"local_pref": 0,
		}
		if r.LocalPref != nil {
			res[i]["local_pref"] = *r.LocalPref
		}
	}

	return res
}
