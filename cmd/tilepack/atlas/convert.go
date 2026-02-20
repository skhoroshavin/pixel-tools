package atlas

import (
	"pixel-tools/cmd/tilepack/tmj"
	"pixel-tools/cmd/tilepack/tsx"
)

func convertTileData(tile *tsx.Tile) map[string]any {
	res := make(map[string]any)

	animation := tmj.ConvertAnimation(tile.Animation)
	if len(animation) > 0 {
		res["animation"] = animation
	}

	objectGroup := tmj.ConvertOptionalObjectGroup(tile.ObjectGroup)
	if objectGroup != nil {
		res["objectgroup"] = objectGroup
	}

	properties := tmj.ConvertProperties(tile.Properties)
	if len(properties) > 0 {
		res["properties"] = properties
	}

	if len(res) == 0 {
		return nil
	}

	return res
}
