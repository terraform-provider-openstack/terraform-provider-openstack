package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandBlockStorageVolumeTypeV3ExtraSpecs(raw map[string]any) volumetypes.ExtraSpecsOpts {
	extraSpecs := make(volumetypes.ExtraSpecsOpts, len(raw))
	for k, v := range raw {
		extraSpecs[k] = v.(string)
	}

	return extraSpecs
}

func expandBlockStorageVolumeTypeV3Encryption(raw []any) volumetypes.CreateEncryptionOpts {
	enc := raw[0].(map[string]any)

	return volumetypes.CreateEncryptionOpts{
		Provider:        enc["provider"].(string),
		Cipher:          enc["cipher"].(string),
		KeySize:         enc["key_size"].(int),
		ControlLocation: enc["control_location"].(string),
	}
}

func flattenBlockStorageVolumeTypeV3Encryption(enc *volumetypes.GetEncryptionType) []map[string]any {
	// Cinder answers 200 with an empty body when the type has no encryption.
	if enc.EncryptionID == "" {
		return nil
	}

	return []map[string]any{
		{
			"provider":         enc.Provider,
			"cipher":           enc.Cipher,
			"key_size":         enc.KeySize,
			"control_location": enc.ControlLocation,
		},
	}
}

// updateBlockStorageVolumeTypeV3Encryption creates, updates or deletes the
// encryption spec to match the configuration. Update and delete need the
// encryption ID, which Cinder returns only from a get.
func updateBlockStorageVolumeTypeV3Encryption(ctx context.Context, client *gophercloud.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("encryption")
	oldEnc, newEnc := oldRaw.([]any), newRaw.([]any)

	if len(oldEnc) == 0 {
		return volumetypes.CreateEncryption(ctx, client, d.Id(), expandBlockStorageVolumeTypeV3Encryption(newEnc)).Err
	}

	current, err := volumetypes.GetEncryption(ctx, client, d.Id()).Extract()
	if err != nil {
		return err
	}

	if len(newEnc) == 0 {
		return volumetypes.DeleteEncryption(ctx, client, d.Id(), current.EncryptionID).ExtractErr()
	}

	updateOpts := volumetypes.UpdateEncryptionOpts(expandBlockStorageVolumeTypeV3Encryption(newEnc))

	return volumetypes.UpdateEncryption(ctx, client, d.Id(), current.EncryptionID, updateOpts).Err
}
