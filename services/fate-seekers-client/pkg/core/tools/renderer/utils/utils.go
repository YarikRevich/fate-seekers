package utils

import (
	"fmt"

	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/tile"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/google/uuid"
	"github.com/lafriks/go-tiled"
)

// LoadMap loads map to different renderer levels.
func LoadMap(tilemap *tiled.Map) {
	for _, layer := range tilemap.Layers {
		layerTiles := loader.GetMapLayerTiles(
			layer, tilemap.Height, tilemap.Width, tilemap.TileHeight, tilemap.TileWidth)

		for _, layerTile := range layerTiles {
			name := uuid.New().String()

			switch layer.Name {
			case loader.FirstMapThirdLayer:
				if !renderer.GetInstance().TertiaryTilemapObjectExists(name) {
					fmt.Println(layerTile.Position)

					renderer.GetInstance().AddTertiaryTilemapObject(
						name, tile.NewTile(layerTile))
				}
			case loader.FirstMapSecondLayer:
			}
		}
	}
}
