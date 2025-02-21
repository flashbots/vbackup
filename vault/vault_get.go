package vault

import "context"

func (v *Vault) Get(ctx context.Context, path string) (interface{}, error) {
	return v.get(ctx, sanitised(path), "")
}

func (v *Vault) get(ctx context.Context, path string, breadcrumbs string) (interface{}, error) {
	// TODO: refactor to v.kvv2?

	{ // if it's a leaf, just get the data and yield
		s, err := v.logical.ReadWithContext(ctx, v.mount+"/data/"+path)
		if err != nil {
			return nil, err
		}
		if s != nil && s.Data != nil {
			_data, ok := s.Data["data"]
			if !ok {
				return nil, nil
			}
			data, ok := _data.(map[string]interface{})
			if !ok {
				return nil, nil
			}
			return data, nil
		}
	}

	{ // otherwise recursively traverse the sub-paths
		s, err := v.logical.ListWithContext(ctx, v.mount+"/metadata/"+path)
		if err != nil {
			return nil, err
		}

		if s == nil || s.Data == nil {
			return nil, nil
		}
		_keys, ok := s.Data["keys"]
		if !ok {
			return nil, nil
		}
		keys, ok := _keys.([]interface{})
		if !ok {
			return nil, nil
		}

		data := make(map[string]interface{}, len(keys))
		for _, _key := range keys {
			key, ok := _key.(string)
			if !ok {
				continue
			}

			key = sanitised(key)
			if _, ignore := v.ignore[sanitised(breadcrumbs+"/"+key)]; ignore {
				continue
			}
			nest, err := v.get(ctx, path+"/"+key, breadcrumbs+"/"+key)
			if err != nil {
				return nil, err
			}
			if nest != nil {
				data[key] = nest
			}
		}

		return data, err
	}
}
